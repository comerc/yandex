Код выдаёт "YES" в отличии от Python и JavaScript

```go
num := 0.1 + 0.1 + 0.1
if num == 0.3 {
  print("YES")
} else {
  print("NO")
}
```

Это связано с тем, как числа с плавающей запятой представлены и обрабатываются в разных языках программирования.

В Python и JavaScript числа с плавающей запятой представлены в формате двоичной системы счисления, что может привести к неточностям при представлении некоторых десятичных дробей. Например, 0.1 в двоичной системе счисления является бесконечной дробью, поэтому она не может быть точно представлена. Это приводит к тому, что сумма 0.1 + 0.1 + 0.1 не равна 0.3.

В Go, с другой стороны, используется библиотека для точной работы с числами с плавающей запятой, которая может обеспечить более точные результаты для таких операций. Ваш код на Go дает “YES”, потому что Go обрабатывает эту ситуацию иначе, чем Python и JavaScript.