from __future__ import annotations

from dataclasses import dataclass
from typing import Any, Dict, List, Tuple

from reportlab.lib.pagesizes import A4
from reportlab.lib.units import mm
from reportlab.pdfbase import pdfmetrics
from reportlab.pdfbase.ttfonts import TTFont
from reportlab.pdfgen import canvas


# Кириллица: положи рядом DejaVuSans.ttf
pdfmetrics.registerFont(TTFont("DejaVuSans", "./DejaVuSans.ttf"))
FONT_NAME = "DejaVuSans"


@dataclass(frozen=True)
class Customer:
    first_name: str
    last_name: str
    phone_number: str


@dataclass(frozen=True)
class VehicleType:
    brand: str
    model: str


@dataclass(frozen=True)
class Vehicle:
    gov_number: str
    vin: str
    mileage: int
    type: VehicleType


@dataclass(frozen=True)
class DocumentData:
    document_id: str
    generated_at: str
    company_name: str
    title: str
    customer: Customer
    vehicle: Vehicle


CHECKLIST_ROWS: List[Tuple[str, str]] = [
    ("Ходовая часть (передняя часть автомобиля)", "Амортизатор (Визуальный осмотр)"),
    ("", "Пружина"),
    ("", "Опора и подшипник амортизатора"),
    ("", "Пыльник, отбойник амортизатора"),
    ("", "Задний, передний сайлент-блок нижнего рычага - подвеска McPherson"),
    ("", "Шаровая опора– подвеска McPherson"),
    ("", "Тяга стабилизатора"),
    ("", "Втулка стабилизатора"),
    ("", "Верхний передний рычаг (сайлент-блок/шаровая опора) - многорычажная подвеска"),
    ("", "Верхний задний рычаг (сайлент-блок/шаровая опора) - многорычажная подвеска"),
    ("", "Нижний передний рычаг (сайлент-блок/шаровая опора) - многорычажная подвеска"),
    ("", "Нижний задний рычаг (сайлент-блок/шаровая опора) - многорычажная подвеска"),
    ("", "Рулевая тяга"),
    ("", "Наконечник рулевой тяги"),
    ("", "Пыльник рулевой тяги"),
    ("", "Люфт, стук, течь рулевой рейки"),
    ("", "Пыльник ШРУСа внешний"),
    ("", "Пыльник ШРУСа внутренний"),
    ("", "Сальник КПП, АКПП"),
    ("", "Ступичный подшипник, ступица"),
    ("Ходовая часть (задняя часть автомобиля)", "Амортизатор"),
    ("", "Пыльник, отбойник амортизатора"),
    ("", "Пружина"),
    ("", "Сайлент-блок задней балки"),
    ("", "Тяга стабилизатора"),
    ("", "Втулки стабилизатора"),
    ("", "Сайлент-блок продольного рычага - многорычажная подвеска"),
    ("", "Нижний несущий рычаг подвески (схожд.) - многорычажная подвеска"),
    ("", "Нижний направляющий рычаг подвески - многорычажная подвеска"),
    ("", "Верхний поперечный рычаг подвески (развал.) - многорычажная подвеска"),
    ("", "Пыльник задней полуоси внутренний – полный привод"),
    ("", "Пыльник задней полуоси внешний– полный привод"),
    ("", "Подшипник ступицы, ступица"),
    ("", "Сальники редуктора– полный привод"),
    ("Прочее", ""),
]


def parse_document(payload: Dict[str, Any]) -> DocumentData:
    cst = payload["customer"]
    v = payload["vehicle"]
    vt = v["type"]

    return DocumentData(
        document_id=payload["document_id"],
        generated_at=payload["generated_at"],
        company_name=payload.get("company_name", "") or "",
        title=payload.get("title", "") or "Документ",
        customer=Customer(
            first_name=cst.get("first_name", "") or "",
            last_name=cst.get("last_name", "") or "",
            phone_number=cst.get("phone_number", "") or "",
        ),
        vehicle=Vehicle(
            gov_number=v.get("gov_number", "") or "",
            vin=v.get("vin", "") or "",
            mileage=int(v.get("mileage") or 0),
            type=VehicleType(
                brand=vt.get("brand", "") or "",
                model=vt.get("model", "") or "",
            ),
        ),
    )


def wrap_text_lines(text: str, font_name: str, font_size: int, max_width: float) -> List[str]:
    """Перенос по словам. Возвращает список строк, каждая <= max_width."""
    words = (text or "").split()
    if not words:
        return [""]

    lines: List[str] = []
    cur = words[0]

    for w in words[1:]:
        candidate = cur + " " + w
        if pdfmetrics.stringWidth(candidate, font_name, font_size) <= max_width:
            cur = candidate
        else:
            lines.append(cur)
            cur = w

    lines.append(cur)
    return lines


def hard_clip_to_width(text: str, font_name: str, font_size: int, max_width: float) -> str:
    """Жёстко обрезает строку до max_width (без троеточий)."""
    if pdfmetrics.stringWidth(text, font_name, font_size) <= max_width:
        return text
    s = text
    while s and pdfmetrics.stringWidth(s, font_name, font_size) > max_width:
        s = s[:-1]
    return s.rstrip()


def count_table_rows(rows: List[Tuple[str, str]]) -> int:
    """
    Сколько фактических строк будет в таблице:
    - 1 строка заголовка таблицы
    - + строки секций
    - + строки элементов
    """
    n = 1  # header row: "Наименование | Лев. | Прав."
    current_section = ""
    for section, item in rows:
        if section and section != current_section:
            current_section = section
            n += 1
        if item:
            n += 1
    return n


def generate_pdf_bytes(payload: Dict[str, Any]) -> bytes:
    data = parse_document(payload)

    from io import BytesIO
    buf = BytesIO()

    c = canvas.Canvas(buf, pagesize=A4)
    w, h = A4

    # Поля страницы
    top_margin = 12 * mm
    bottom_margin = 8 * mm

    # Заголовок
    c.setFont(FONT_NAME, 13)
    c.drawString(20 * mm, h - top_margin, data.title or "Карта проверки и осмотра ходовой части автомобиля")

    # Шапка
    header_font = 9
    c.setFont(FONT_NAME, header_font)

    y = h - top_margin - 8 * mm
    line_h = 5 * mm

    brand = data.vehicle.type.brand
    model = data.vehicle.type.model
    customer_name = (data.customer.first_name + " " + data.customer.last_name).strip()

    header_lines = [
        ("марка", brand),
        ("модель", model),
        ("гос. номер", data.vehicle.gov_number),
        ("VIN", data.vehicle.vin),
        ("пробег", str(data.vehicle.mileage)),
        ("заказчик", customer_name),
        ("контактный номер заказчика", data.customer.phone_number),
    ]

    label_x = 20 * mm
    max_label_w = max(pdfmetrics.stringWidth(f"{lbl}:", FONT_NAME, header_font) for lbl, _ in header_lines)
    value_x = label_x + max_label_w + 6 * mm

    for lbl, val in header_lines:
        c.drawString(label_x, y, f"{lbl}:")
        c.drawString(value_x, y, val)
        y -= line_h

    y -= 2 * mm  # небольшой отступ до таблицы

    # Таблица: размеры колонок
    left_x = 20 * mm
    name_w = 120 * mm
    col_w = 25 * mm

    # Рассчитать высоту строки так, чтобы таблица влезла на 1 лист
    total_rows = count_table_rows(CHECKLIST_ROWS)

    available_h = y - bottom_margin - 4 * mm  # чуть воздуха снизу
    # Предпочитаем 6мм, но если не влезает — уменьшаем
    row_h = min(6 * mm, max(3.8 * mm, available_h / total_rows))

    # Шрифты для таблицы: тоже подстраиваем
    # (при маленьких строках уменьшаем кегль, чтобы текст помещался)
    if row_h >= 5.8 * mm:
        section_font = 8
        item_font = 7
    elif row_h >= 4.8 * mm:
        section_font = 7
        item_font = 6
    else:
        section_font = 6
        item_font = 5

    padding_x = 2 * mm
    padding_top = 1.2 * mm

    def draw_row_box(y0: float) -> None:
        c.rect(left_x, y0 - row_h, name_w, row_h, stroke=1, fill=0)
        c.rect(left_x + name_w, y0 - row_h, col_w, row_h, stroke=1, fill=0)
        c.rect(left_x + name_w + col_w, y0 - row_h, col_w, row_h, stroke=1, fill=0)

    def draw_name_cell(y0: float, text: str, is_section: bool) -> None:
        font_size = section_font if is_section else item_font
        c.setFont(FONT_NAME, font_size)

        max_text_w = name_w - 2 * padding_x

        # Пытаемся в 2 строки (если по высоте позволяет)
        # 2 строки реально читаемы если row_h >= ~5мм
        allow_two_lines = row_h >= 5.0 * mm

        lines = wrap_text_lines(text, FONT_NAME, font_size, max_text_w)
        if allow_two_lines and len(lines) >= 2:
            l1 = lines[0]
            l2 = " ".join(lines[1:])  # всё остальное во 2 строку
            l2 = hard_clip_to_width(l2, FONT_NAME, font_size, max_text_w)

            # Вертикальные позиции строк внутри ячейки
            # первая строка ближе к верху, вторая ниже
            y_text1 = y0 - padding_top - font_size * 0.2
            y_text2 = y_text1 - (font_size + 1)  # небольшой интерлиньяж

            c.drawString(left_x + padding_x, y_text1 - row_h + (row_h - (font_size + 1)) , l1)  # страховка
            # Лучше проще: фиксируем от верхней границы ячейки
            top_of_cell = y0
            c.drawString(left_x + padding_x, top_of_cell - padding_top - font_size, l1)
            c.drawString(left_x + padding_x, top_of_cell - padding_top - font_size - (font_size + 1), l2)
            return

        # 1 строка (или row_h слишком мал) — переносим в 1 строку и жёстко режем
        one = " ".join(lines)
        one = hard_clip_to_width(one, FONT_NAME, font_size, max_text_w)
        c.drawString(left_x + padding_x, y0 - row_h + 1.3 * mm, one)

    def draw_row(y0: float, name: str, is_section: bool = False) -> None:
        draw_row_box(y0)
        draw_name_cell(y0, name, is_section=is_section)

    # Header row таблицы
    c.setFont(FONT_NAME, section_font)
    draw_row_box(y)
    c.drawString(left_x + padding_x, y - row_h + 1.3 * mm, "Наименование")
    c.drawString(left_x + name_w + 7 * mm, y - row_h + 1.3 * mm, "Лев.")
    c.drawString(left_x + name_w + col_w + 6 * mm, y - row_h + 1.3 * mm, "Прав.")
    y -= row_h

    # Rows
    current_section = ""
    for section, item in CHECKLIST_ROWS:
        if section and section != current_section:
            current_section = section
            draw_row(y, section, is_section=True)
            y -= row_h

        if item:
            draw_row(y, item, is_section=False)
            y -= row_h

    # Низ: document_id / generated_at (очень мелко, чтобы не мешало 1 странице)
    c.setFont(FONT_NAME, 6)
    footer = f"document_id={data.document_id} generated_at={data.generated_at}"
    c.drawString(20 * mm, 4 * mm, footer)

    c.showPage()
    c.save()
    return buf.getvalue()


if __name__ == "__main__":
    import json
    from pathlib import Path

    payload = json.loads(Path("test.json").read_text(encoding="utf-8"))
    pdf = generate_pdf_bytes(payload)
    Path("out.pdf").write_bytes(pdf)
    print("Wrote out.pdf")
