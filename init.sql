-- Создаем базу данных Dota
CREATE DATABASE Dota;

-- Подключаемся к базе данных Dota и создаем таблицы из дампа
\connect Dota;
-- Выполните скрипт dota_dump.sql
\i /docker-entrypoint-initdb.d/Dota.dump;

-- Создаем базу данных DotaTest
CREATE DATABASE DotaTest;

-- Подключаемся к базе данных DotaTest и создаем таблицы из дампа
\connect DotaTest;
-- Выполните скрипт dota_test_dump.sql
\i /docker-entrypoint-initdb.d/DotaTest.dump;
