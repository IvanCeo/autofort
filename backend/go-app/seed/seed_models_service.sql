-- Полная матрица услуг для первых 5 автомобилей (все 50 услуг для каждого)
-- Lada Granta: 'a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c'
INSERT INTO vehicle_type_services (vehicle_type_id, service_id, estimated_time_minutes, created_at) VALUES
-- Основные ТО
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8a3f7b2c-4d9e-1a6b-5c2f-8e4d3a7b1c9f', 45, now()),  -- Замена моторного масла
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '3b8c5d9a-2e7f-4a1b-6c3d-9f5e2b8c4a1d', 15, now()),  -- Замена воздушного фильтра
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '7c2a9d4b-5e3f-8b1c-6a9d-2f4c7e5a3b8d', 20, now()),  -- Замена салонного фильтра
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '4d6b8c2a-9e5f-3a7b-1c8d-5f2e4b9d6a3c', 50, now()),  -- Замена топливного фильтра
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '9e3c7a5b-2d8f-4b6c-1a9e-3f5d2c8b4a7e', 30, now()),  -- Замена свечей зажигания
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '2a7d4c9b-6e3f-5b8a-1c9d-4f2e7a5b3c8d', 180, now()), -- Замена ремня ГРМ
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '5c9e3b8a-7d2f-4a6c-1b9e-8f3d5c2a7b4e', 240, now()), -- Замена цепи ГРМ
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8b2d7c4a-9e6f-3a5b-1c8d-4f7e2b9d5a3c', 40, now()),  -- Замена приводных ремней

-- Тормозная система
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '6a4c9d3b-8e2f-5b7a-1c9d-3f6e4a8b2c5d', 60, now()),  -- Замена тормозных колодок
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '3d7b5c9a-2e8f-4a6b-1c3d-9f5e7b2c4a8d', 120, now()), -- Замена тормозных дисков
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '9b4e7c2a-5d8f-3a6b-1c9e-4f2d5b8a3c7e', 100, now()), -- Замена тормозных барабанов
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '7a2d9c4b-6e3f-5b8a-1c7d-8f4e2a5b3c9d', 50, now()),  -- Прокачка тормозной системы

-- Подвеска и ходовая
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '4c8b3d9a-7e2f-5a6b-1c4d-9f3e8b2c5a7d', 120, now()), -- Замена амортизаторов
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '2b9d5c7a-4e3f-8a6b-1c2d-5f7e9b3c4a8d', 40, now()),  -- Замена стоек стабилизатора
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '6d3a8c5b-9e2f-4a7b-1c6d-3f5e9a2c8b4d', 90, now()),  -- Замена шаровых опор
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8c5b2d9a-7e4f-3a6b-1c8d-5f2e7c4b3a9d', 120, now()), -- Замена сайлентблоков
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '5a7c4b9d-2e8f-3a6b-1c5d-9f4e7a2c8b3d', 80, now()),  -- Замена опорных подшипников
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '3b6d8c5a-9e2f-4a7b-1c3d-8f5e6b2c9a4d', 60, now()),  -- Замена рулевых наконечников

-- Рулевое управление
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '9d2b7c4a-6e3f-5a8b-1c9d-4f2e8b3c5a7d', 240, now()), -- Замена рулевой рейки
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '7c4a9d3b-8e2f-5b6a-1c7d-3f9e4c2a8b5d', 120, now()), -- Замена насоса ГУР
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '4a8c5b3d-9e2f-7a6b-1c4d-8f3e5a2c9b7d', 30, now()),  -- Замена жидкости ГУР
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '6b3d8c5a-7e4f-2a9b-1c6d-5f3e8b4c2a7d', 60, now()),  -- Регулировка развала-схождения

-- Трансмиссия
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '2c9b5d7a-4e8f-3a6b-1c2d-9f5e7c4b3a8d', 90, now()),  -- Замена масла в АКПП
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8d3a6c5b-9e7f-4a2b-1c8d-3f5e9a4c7b2d', 60, now()),  -- Замена масла в МКПП
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '5b7c3d9a-2e8f-4a6b-1c5d-9f3e7b2c4a8d', 240, now()), -- Замена сцепления
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '3c8b5d7a-9e4f-2a6b-1c3d-8f5e7c4b2a9d', 120, now()), -- Замена ШРУСа
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '9a4d7c2b-6e3f-5a8b-1c9d-4f2e7a5c3b8d', 60, now()),  -- Замена пыльника ШРУСа

-- Двигатель и системы
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '7d2b9c4a-5e8f-3a6b-1c7d-9f4e2b5c8a3d', 150, now()), -- Замена радиатора охлаждения
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '4c7a8d3b-9e5f-2a6b-1c4d-8f3e7c5a2b9d', 45, now()),  -- Замена термостата
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '6a3d8c5b-7e9f-4a2b-1c6d-5f3e8a4c7b2d', 90, now()),  -- Замена помпы
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8b5c3d7a-2e9f-4a6b-1c8d-5f3e7b4c2a8d', 210, now()), -- Замена ремня ГРМ с помпой
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '3a6c8d5b-9e7f-2a4b-1c3d-8f5e6a4c7b2d', 120, now()), -- Чистка инжектора
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '5d7b3c9a-8e4f-2a6b-1c5d-9f3e7b4c2a8d', 40, now()),  -- Чистка дроссельной заслонки
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '2b8c6d3a-7e9f-4a5b-1c2d-8f3e7b4c5a9d', 50, now()),  -- Замена катушек зажигания

-- Электрика и аккумулятор
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '9c4a7d2b-6e8f-3a5b-1c9d-4f2e7a5c3b8d', 20, now()),  -- Замена аккумулятора
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '7a3d9c5b-8e4f-2a6b-1c7d-5f3e9a4c7b2d', 120, now()), -- Замена генератора
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '4b8c5d3a-9e7f-2a6b-1c4d-8f3e5b4c7a2d', 90, now()),  -- Замена стартера
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '6d2a8c5b-7e9f-3a4b-1c6d-5f3e8a4c7b2d', 30, now()),  -- Замена ламп освещения
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8c5b3d7a-2e9f-4a6b-1c8d-5f3e7b4c2a8d', 15, now()),  -- Замена дворников

-- Кондиционер и отопление
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '3b6c8d5a-9e7f-2a4b-1c3d-8f5e6a4c7b2d', 60, now()),  -- Заправка кондиционера
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '5a7c4d3b-8e9f-2a6b-1c5d-4f3e7a5c2b8d', 180, now()), -- Замена компрессора кондиционера
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '2c9b6d3a-7e8f-4a5b-1c2d-9f3e7b4c5a8d', 150, now()), -- Замена радиатора кондиционера
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8d3a7c5b-9e6f-2a4b-1c8d-5f3e8a4c7b2d', 120, now()), -- Замена печки

-- Шины и колеса
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '6a4d8c3b-7e9f-2a5b-1c6d-8f3e5a4c7b2d', 90, now()),  -- Сезонная замена шин
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '4b7c5d3a-9e8f-2a6b-1c4d-5f3e7a4c2b9d', 60, now()),  -- Балансировка колес
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '9d2a7c4b-6e8f-3a5b-1c9d-5f4e2a7c3b8d', 20, now()),  -- Ремонт прокола шины
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '7c3a9d5b-8e6f-2a4b-1c7d-5f3e9a4c7b2d', 10, now()),  -- Замена вентилей

-- Кузовные работы
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '5b8c4d3a-9e7f-2a6b-1c5d-8f3e5b4c7a2d', 180, now()), -- Замена лобового стекла
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '3a7c6d5b-8e9f-2a4b-1c3d-5f3e7a4c2b8d', 90, now()),  -- Полировка фар
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8d2b7c4a-6e9f-3a5b-1c8d-5f4e2a7c3b9d', 240, now()), -- Химчистка салона
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '6c3a8d5b-7e9f-2a4b-1c6d-5f3e8a4c7b2d', 300, now()), -- Покрытие керамикой

-- Диагностика
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '4a9c6d3b-8e7f-2a5b-1c4d-5f3e7a4c2b8d', 30, now()),  -- Компьютерная диагностика
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '2b7c8d3a-9e6f-4a5b-1c2d-8f3e7b4c5a9d', 45, now()),  -- Диагностика подвески
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '9c3a7d5b-8e6f-2a4b-1c9d-5f3e8a4c7b2d', 20, now()),  -- Считывание ошибок
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '7d2a9c4b-6e8f-3a5b-1c7d-5f4e2a7c3b8d', 25, now()),  -- Адаптация дроссельной заслонки

-- Дополнительные услуги
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '5a8c7d3b-9e6f-2a4b-1c5d-8f3e5a4c7b2d', 90, now()),  -- Замена глушителя
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '3b6c9d5a-8e7f-2a4b-1c3d-5f3e7a4c2b8d', 150, now()), -- Замена сажевого фильтра
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8c4a7d3b-9e6f-2a5b-1c8d-5f3e8a4c7b2d', 120, now()), -- Замена катализатора
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '6d3a9c5b-7e8f-2a4b-1c6d-5f3e8a4c7b2d', 60, now()),  -- Установка защиты картера
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '4b8c6d3a-9e7f-2a5b-1c4d-8f3e5b4c7a2d', 180, now()); -- Антикоррозийная обработка

INSERT INTO vehicle_type_services (vehicle_type_id, service_id, estimated_time_minutes, created_at) VALUES
-- Lada Granta
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8a3f7b2c-4d9e-1a6b-5c2f-8e4d3a7b1c9f', 45, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '3b8c5d9a-2e7f-4a1b-6c3d-9f5e2b8c4a1d', 15, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '7c2a9d4b-5e3f-8b1c-6a9d-2f4c7e5a3b8d', 20, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '4d6b8c2a-9e5f-3a7b-1c8d-5f2e4b9d6a3c', 50, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '9e3c7a5b-2d8f-4b6c-1a9e-3f5d2c8b4a7e', 30, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '2a7d4c9b-6e3f-5b8a-1c9d-4f2e7a5b3c8d', 180, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '5c9e3b8a-7d2f-4a6c-1b9e-8f3d5c2a7b4e', 240, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8b2d7c4a-9e6f-3a5b-1c8d-4f7e2b9d5a3c', 40, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '6a4c9d3b-8e2f-5b7a-1c9d-3f6e4a8b2c5d', 60, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '3d7b5c9a-2e8f-4a6b-1c3d-9f5e7b2c4a8d', 120, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '9b4e7c2a-5d8f-3a6b-1c9e-4f2d5b8a3c7e', 100, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '7a2d9c4b-6e3f-5b8a-1c7d-8f4e2a5b3c9d', 50, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '4c8b3d9a-7e2f-5a6b-1c4d-9f3e8b2c5a7d', 120, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '2b9d5c7a-4e3f-8a6b-1c2d-5f7e9b3c4a8d', 40, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '6d3a8c5b-9e2f-4a7b-1c6d-3f5e9a2c8b4d', 90, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8c5b2d9a-7e4f-3a6b-1c8d-5f2e7c4b3a9d', 120, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '5a7c4b9d-2e8f-3a6b-1c5d-9f4e7a2c8b3d', 80, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '3b6d8c5a-9e2f-4a7b-1c3d-8f5e6b2c9a4d', 60, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '9d2b7c4a-6e3f-5a8b-1c9d-4f2e8b3c5a7d', 240, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '7c4a9d3b-8e2f-5b6a-1c7d-3f9e4c2a8b5d', 120, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '4a8c5b3d-9e2f-7a6b-1c4d-8f3e5a2c9b7d', 30, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '6b3d8c5a-7e4f-2a9b-1c6d-5f3e8b4c2a7d', 60, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '2c9b5d7a-4e8f-3a6b-1c2d-9f5e7c4b3a8d', 90, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8d3a6c5b-9e7f-4a2b-1c8d-3f5e9a4c7b2d', 60, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '5b7c3d9a-2e8f-4a6b-1c5d-9f3e7b2c4a8d', 240, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '3c8b5d7a-9e4f-2a6b-1c3d-8f5e7c4b2a9d', 120, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '9a4d7c2b-6e3f-5a8b-1c9d-4f2e7a5c3b8d', 60, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '7d2b9c4a-5e8f-3a6b-1c7d-9f4e2b5c8a3d', 150, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '4c7a8d3b-9e5f-2a6b-1c4d-8f3e7c5a2b9d', 45, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '6a3d8c5b-7e9f-4a2b-1c6d-5f3e8a4c7b2d', 90, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8b5c3d7a-2e9f-4a6b-1c8d-5f3e7b4c2a8d', 210, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '3a6c8d5b-9e7f-2a4b-1c3d-8f5e6a4c7b2d', 120, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '5d7b3c9a-8e4f-2a6b-1c5d-9f3e7b4c2a8d', 40, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '2b8c6d3a-7e9f-4a5b-1c2d-8f3e7b4c5a9d', 50, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '9c4a7d2b-6e8f-3a5b-1c9d-4f2e7a5c3b8d', 20, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '7a3d9c5b-8e4f-2a6b-1c7d-5f3e9a4c7b2d', 120, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '4b8c5d3a-9e7f-2a6b-1c4d-8f3e5b4c7a2d', 90, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '6d2a8c5b-7e9f-3a4b-1c6d-5f3e8a4c7b2d', 30, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8c5b3d7a-2e9f-4a6b-1c8d-5f3e7b4c2a8d', 15, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '3b6c8d5a-9e7f-2a4b-1c3d-8f5e6a4c7b2d', 60, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '5a7c4d3b-8e9f-2a6b-1c5d-4f3e7a5c2b8d', 180, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '2c9b6d3a-7e8f-4a5b-1c2d-9f3e7b4c5a8d', 150, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8d3a7c5b-9e6f-2a4b-1c8d-5f3e8a4c7b2d', 120, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '6a4d8c3b-7e9f-2a5b-1c6d-8f3e5a4c7b2d', 90, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '4b7c5d3a-9e8f-2a6b-1c4d-5f3e7a4c2b9d', 60, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '9d2a7c4b-6e8f-3a5b-1c9d-5f4e2a7c3b8d', 20, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '7c3a9d5b-8e6f-2a4b-1c7d-5f3e9a4c7b2d', 10, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '5a8c7d3b-9e6f-2a4b-1c5d-8f3e5a4c7b2d', 90, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '3b6c9d5a-8e7f-2a4b-1c3d-5f3e7a4c2b8d', 150, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '8c4a7d3b-9e6f-2a5b-1c8d-5f3e8a4c7b2d', 120, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '6d3a9c5b-7e8f-2a4b-1c6d-5f3e8a4c7b2d', 60, now()),
('a8e1f4b2-7c3d-4f9a-8b6c-5d2e1f3a4b7c', '4b8c6d3a-9e7f-2a5b-1c4d-8f3e5b4c7a2d', 180, now()),

-- Lada Vesta (время на 10% больше чем у Granta)
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '8a3f7b2c-4d9e-1a6b-5c2f-8e4d3a7b1c9f', 50, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '3b8c5d9a-2e7f-4a1b-6c3d-9f5e2b8c4a1d', 17, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '7c2a9d4b-5e3f-8b1c-6a9d-2f4c7e5a3b8d', 22, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '4d6b8c2a-9e5f-3a7b-1c8d-5f2e4b9d6a3c', 55, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '9e3c7a5b-2d8f-4b6c-1a9e-3f5d2c8b4a7e', 33, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '2a7d4c9b-6e3f-5b8a-1c9d-4f2e7a5b3c8d', 198, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '5c9e3b8a-7d2f-4a6c-1b9e-8f3d5c2a7b4e', 264, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '8b2d7c4a-9e6f-3a5b-1c8d-4f7e2b9d5a3c', 44, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '6a4c9d3b-8e2f-5b7a-1c9d-3f6e4a8b2c5d', 66, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '3d7b5c9a-2e8f-4a6b-1c3d-9f5e7b2c4a8d', 132, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '9b4e7c2a-5d8f-3a6b-1c9e-4f2d5b8a3c7e', 110, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '7a2d9c4b-6e3f-5b8a-1c7d-8f4e2a5b3c9d', 55, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '4c8b3d9a-7e2f-5a6b-1c4d-9f3e8b2c5a7d', 132, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '2b9d5c7a-4e3f-8a6b-1c2d-5f7e9b3c4a8d', 44, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '6d3a8c5b-9e2f-4a7b-1c6d-3f5e9a2c8b4d', 99, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '8c5b2d9a-7e4f-3a6b-1c8d-5f2e7c4b3a9d', 132, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '5a7c4b9d-2e8f-3a6b-1c5d-9f4e7a2c8b3d', 88, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '3b6d8c5a-9e2f-4a7b-1c3d-8f5e6b2c9a4d', 66, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '9d2b7c4a-6e3f-5a8b-1c9d-4f2e8b3c5a7d', 264, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '7c4a9d3b-8e2f-5b6a-1c7d-3f9e4c2a8b5d', 132, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '4a8c5b3d-9e2f-7a6b-1c4d-8f3e5a2c9b7d', 33, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '6b3d8c5a-7e4f-2a9b-1c6d-5f3e8b4c2a7d', 66, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '2c9b5d7a-4e8f-3a6b-1c2d-9f5e7c4b3a8d', 99, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '8d3a6c5b-9e7f-4a2b-1c8d-3f5e9a4c7b2d', 66, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '5b7c3d9a-2e8f-4a6b-1c5d-9f3e7b2c4a8d', 264, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '3c8b5d7a-9e4f-2a6b-1c3d-8f5e7c4b2a9d', 132, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '9a4d7c2b-6e3f-5a8b-1c9d-4f2e7a5c3b8d', 66, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '7d2b9c4a-5e8f-3a6b-1c7d-9f4e2b5c8a3d', 165, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '4c7a8d3b-9e5f-2a6b-1c4d-8f3e7c5a2b9d', 50, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '6a3d8c5b-7e9f-4a2b-1c6d-5f3e8a4c7b2d', 99, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '8b5c3d7a-2e9f-4a6b-1c8d-5f3e7b4c2a8d', 231, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '3a6c8d5b-9e7f-2a4b-1c3d-8f5e6a4c7b2d', 132, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '5d7b3c9a-8e4f-2a6b-1c5d-9f3e7b4c2a8d', 44, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '2b8c6d3a-7e9f-4a5b-1c2d-8f3e7b4c5a9d', 55, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '9c4a7d2b-6e8f-3a5b-1c9d-4f2e7a5c3b8d', 22, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '7a3d9c5b-8e4f-2a6b-1c7d-5f3e9a4c7b2d', 132, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '4b8c5d3a-9e7f-2a6b-1c4d-8f3e5b4c7a2d', 99, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '6d2a8c5b-7e9f-3a4b-1c6d-5f3e8a4c7b2d', 33, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '8c5b3d7a-2e9f-4a6b-1c8d-5f3e7b4c2a8d', 17, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '3b6c8d5a-9e7f-2a4b-1c3d-8f5e6a4c7b2d', 66, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '5a7c4d3b-8e9f-2a6b-1c5d-4f3e7a5c2b8d', 198, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '2c9b6d3a-7e8f-4a5b-1c2d-9f3e7b4c5a8d', 165, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '8d3a7c5b-9e6f-2a4b-1c8d-5f3e8a4c7b2d', 132, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '6a4d8c3b-7e9f-2a5b-1c6d-8f3e5a4c7b2d', 99, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '4b7c5d3a-9e8f-2a6b-1c4d-5f3e7a4c2b9d', 66, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '9d2a7c4b-6e8f-3a5b-1c9d-5f4e2a7c3b8d', 22, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '7c3a9d5b-8e6f-2a4b-1c7d-5f3e9a4c7b2d', 11, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '5a8c7d3b-9e6f-2a4b-1c5d-8f3e5a4c7b2d', 99, now()),
('b5c9d2a1-8e4f-3b7c-9d6a-1f2e4c5b8a3d', '3