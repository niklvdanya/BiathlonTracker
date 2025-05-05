# BiathlonTracker
Тестовое задание Yadro Impulse 2025

## Описание
Система отслеживания соревнований по биатлону с возможностью обработки различных событий, расчета времени прохождения кругов и формирования итоговых результатов.

## Сборка
```bash
make build
```
Или вручную:
```bash
go build -o biathlon ./cmd/main.go
```

## Запуск
```bash
./biathlon
```
С дополнительными параметрами:
```bash
./biathlon -config=config.json -events=events.txt -parallel
```

## Тесты
Запуск всех тестов:
```bash
make test
```