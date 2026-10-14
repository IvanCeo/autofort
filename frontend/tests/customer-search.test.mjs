import assert from "node:assert/strict"
import { readFile } from "node:fs/promises"
import test from "node:test"
import { compileScript, parse } from "@vue/compiler-sfc"
import * as Vue from "vue"

// Compile the real SFCs using the Vue tooling already installed with Nuxt.
// Supply Nuxt auto-imports and deterministic requests without a browser/server.
async function loadComponent(path, autoImports = {}) {
  const source = await readFile(new URL(path, import.meta.url), "utf8")
  const { descriptor } = parse(source)
  const { content } = compileScript(descriptor, {
    id: path,
    inlineTemplate: true,
    genDefaultAs: "component",
  })
  const executable = content.replace(
    /import\s*\{([^}]+)\}\s*from\s*["']vue["'];?/g,
    (_, imports) => `const { ${imports.replace(/\s+as\s+/g, ": ")} } = Vue;`
  )
  const bindings = { ref: Vue.ref, computed: Vue.computed, watch: Vue.watch, ...autoImports }
  return new Function("Vue", ...Object.keys(bindings), `${executable}\nreturn component`)(
    Vue,
    ...Object.values(bindings)
  )
}

function deferred() {
  let resolve, reject
  const promise = new Promise((yes, no) => { resolve = yes; reject = no })
  return {
    promise,
    settled: false,
    resolve(value) { this.settled = true; resolve(value) },
    reject(error) { this.settled = true; reject(error) },
  }
}

// A small host tree lets Vue mount real components, directives and watchers.
// It supplies input listeners but does not model browser layout or networking.
function createHost() {
  const node = (tag, text = "") => ({ tag, text, props: {}, children: [], parent: null })
  const remove = (el) => {
    if (el.parent) el.parent.children.splice(el.parent.children.indexOf(el), 1)
    el.parent = null
  }
  const insert = (el, parent, anchor = null) => {
    remove(el)
    parent.children.splice(anchor ? parent.children.indexOf(anchor) : parent.children.length, 0, el)
    el.parent = parent
  }
  const renderer = Vue.createRenderer({
    createElement(tag) {
      const el = node(tag)
      el.value = ""
      el.listeners = {}
      el.addEventListener = (event, callback) => { el.listeners[event] = callback }
      return el
    },
    createText: (text) => node("#text", text),
    createComment: () => node("#comment"),
    insert,
    remove,
    setText: (el, text) => { el.text = text },
    setElementText: (el, text) => { el.text = text; el.children = [] },
    patchProp: (el, key, previous, value) => { el.props[key] = value },
    parentNode: (el) => el.parent,
    nextSibling: (el) => el.parent?.children[el.parent.children.indexOf(el) + 1] || null,
    insertStaticContent(content, parent, anchor) {
      const el = node("#static", content)
      insert(el, parent, anchor)
      return [el, el]
    },
  })
  return { renderer, root: node("root") }
}

function find(root, predicate) {
  if (predicate(root)) return root
  for (const child of root.children) {
    const match = find(child, predicate)
    if (match) return match
  }
}

function text(root) {
  return root.text + root.children.map(text).join("")
}

async function flush() {
  await new Promise((resolve) => setImmediate(resolve))
  await Vue.nextTick()
}

async function mountSearch() {
  const requests = []
  const warnings = []
  const timers = new Map()
  const dictionary = deferred()
  const search = deferred()
  const state = { data: Vue.ref(), pending: Vue.ref(false), error: Vue.ref(null) }

  function dictionaryFetch(url, lazy) {
    const start = () => {
      requests.push({ url })
      state.pending.value = true
      return dictionary.promise.then(
        (result) => { state.data.value = result },
        (error) => { state.error.value = error }
      ).finally(() => { state.pending.value = false })
    }
    // Nuxt SPA lazy fetching starts onBeforeMount; awaiting its state does not
    // await the network request. Eager useFetch awaits that request instead.
    if (lazy) {
      Vue.onBeforeMount(start)
      return Object.assign(Promise.resolve(state), state)
    }
    return Object.assign(start().then(() => state), state)
  }

  const page = await loadComponent("../app/pages/index.vue", {
    useFetch: (url) => dictionaryFetch(url, false),
    useLazyFetch: (url) => dictionaryFetch(url, true),
    $fetch: (url, options) => {
      requests.push({ url, ...options })
      return search.promise
    },
    setTimeout: (callback, delay) => {
      const id = Symbol("timer")
      timers.set(id, { callback, delay })
      return id
    },
    clearTimeout: (id) => timers.delete(id),
  })
  const card = await loadComponent("../app/components/CustomerCard.vue")
  const { renderer, root } = createHost()
  const app = renderer.createApp({
    setup: () => () => Vue.h(Vue.Suspense, null, {
      default: () => Vue.h(page),
      fallback: () => Vue.h("p", "Page pending"),
    }),
  })
  app.component("CustomerCard", card)
  app.component("NuxtLink", {
    props: ["to"],
    setup: (props, { slots }) => () => Vue.h("a", { href: props.to }, slots.default?.()),
  })
  app.config.warnHandler = (warning) => warnings.push(warning)
  app.mount(root)
  await flush()
  return { app, root, dictionary, search, timers, requests, warnings }
}

async function searchWhileDictionaryPending(fixture) {
  const { root, dictionary, search, timers, requests } = fixture
  assert.equal(dictionary.settled, false)
  const input = find(root, (el) => el.tag === "input")
  assert.ok(input, "search input must mount before the dictionary settles")
  assert.equal(typeof input.listeners.input, "function")
  input.value = "Иван"
  input.listeners.input({ target: input })
  await flush()
  assert.equal(timers.size, 1, "the real Vue watcher must schedule search")
  const [id, timer] = timers.entries().next().value
  assert.equal(timer.delay, 300)
  timers.delete(id)
  const action = timer.callback()
  await flush()
  assert.match(text(root), /Ищем\.\.\./)
  assert.deepEqual(requests, [
    { url: "/api/vehicle-types" },
    { url: "/api/customers", method: "GET", query: { q: "Иван" } },
  ])
  search.resolve({
    items: [{
      id: "customer-1",
      first_name: "Иван",
      last_name: "Иванов",
      phone_number: "+79991234567",
      vehicles: [{ vehicle_type_id: "vehicle-type-1" }],
    }],
  })
  await action
  await flush()
  assert.equal(dictionary.settled, false, "dictionary must remain unresolved during search")
  assert.match(text(root), /Иван Иванов/)
  assert.match(text(root), /Unknown model/)
  assert.match(text(root), /Загружаем машины/)
  assert.doesNotMatch(text(root), /Ищем\.\.\./)
  assert.ok(find(root, (el) => el.props.href === "/customers/customer-1"))
}

test("home search works before dictionary resolution and reactively displays brand/model", async (t) => {
  const fixture = await mountSearch()
  t.after(() => fixture.app.unmount())
  await searchWhileDictionaryPending(fixture)
  fixture.dictionary.resolve({ items: [{ id: "vehicle-type-1", brand: "Toyota", model: "Camry" }] })
  await flush()
  assert.match(text(fixture.root), /Toyota Camry/)
  assert.doesNotMatch(text(fixture.root), /Unknown model|Загружаем машины/)
  assert.deepEqual(fixture.warnings, [])
})

test("home search keeps existing results when the pending dictionary fails", async (t) => {
  const fixture = await mountSearch()
  t.after(() => fixture.app.unmount())
  await searchWhileDictionaryPending(fixture)
  fixture.dictionary.reject(new Error("Dictionary unavailable"))
  await flush()
  assert.match(text(fixture.root), /Dictionary unavailable/)
  assert.match(text(fixture.root), /Иван Иванов/)
  assert.match(text(fixture.root), /Unknown model/)
  assert.doesNotMatch(text(fixture.root), /Загружаем машины/)
  assert.deepEqual(fixture.warnings, [])
})
