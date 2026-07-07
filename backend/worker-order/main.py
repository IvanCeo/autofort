from __future__ import annotations

import json
from typing import Any, Dict

import anyio
from fastapi import FastAPI, Request
from fastapi.responses import Response

from pdfgen import generate_pdf_bytes

app = FastAPI(title="work-order-pdf-service", version="0.1.0")


@app.post("/work-order/pdf")
async def work_order_pdf(request: Request) -> Response:
    """
    Принимает JSON (как DTO из Go), возвращает application/pdf.
    Асинхронная ручка, генерацию PDF выносим в threadpool.
    """
    try:
        payload: Dict[str, Any] = await request.json()
    except json.JSONDecodeError:
        return Response(content="invalid JSON", status_code=400, media_type="text/plain")
    except Exception as e:
        return Response(content=f"failed to read JSON: {e}", status_code=400, media_type="text/plain")

    # Минимальная валидация (без Pydantic, чтобы не городить):
    required_top = ["document_id", "generated_at", "title", "customer", "vehicle"]
    for k in required_top:
        if k not in payload:
            return Response(content=f"missing field: {k}", status_code=400, media_type="text/plain")

    try:
        pdf_bytes = await anyio.to_thread.run_sync(generate_pdf_bytes, payload)
    except KeyError as e:
        # если внутри parse_document не нашлось поля
        return Response(content=f"missing field in payload: {e}", status_code=400, media_type="text/plain")
    except Exception as e:
        # чтобы Go-клиент получил диагностируемую ошибку (он читает тело)
        return Response(content=f"pdf generation error: {e}", status_code=500, media_type="text/plain")

    # Важно: Content-Type ровно application/pdf
    return Response(
        content=pdf_bytes,
        status_code=200,
        media_type="application/pdf",
        headers={
            # опционально, но удобно
            "Content-Disposition": f'inline; filename="{payload.get("document_id", "work-order")}.pdf"',
        },
    )
