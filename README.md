## Текст задания:

В программе есть слайс, содержащий URL. На каждый URL нужно отправить HTTP-запрос методом GET
и посчитать кол-во вхождений строки "Go" в теле ответа. В конце работы приложение выводит на экран общее количество найденных строк "Go" во всех переданных URL, например:

```
Count for https://golang.org: 20
Count for https://golang.org: 20
Total: 40
```

Каждый URL должен начать обрабатываться сразу после вычитывания и параллельно с вычитыванием следующего. URL должны обрабатываться параллельно, но не более k=5 одновременно. Обработчики URL не должны порождать лишних горутин, т.е. если k=5, а обрабатываемых URL-ов всего 2, не должно создаваться 5 горутин.

Нужно обойтись без глобальных переменных и использовать только стандартную библиотеку.
>>>>>>> 5c5add22449734e81c6f28007ff042e873f33c1d
