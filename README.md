# capital-gains
CLI app para calcular imposto sobre operações de venda e compra em cima de uma ou mais ações

## Começando

### Para rodar esse app CLI:
- Essa aplicação usa apenas a biblioteca padrão, não sendo necessário portanto, a instalação de dependências externas.
```Bash
go run cmd/main.go < input.txt
```

- Opcionalmente podemos passar argumentos para flexibilizar o comportamento da aplicaçãp, ambos são do tipo `float64`.
 
 ***Valores default:***

`PROFIT_PERCENTAGE` = 0.20

`MAX_SELL_OPERATION_VALUE` = 20000.0

```Bash
PROFIT_PERCENTAGE=0.30 MAX_SELL_OPERATION_VALUE=25000.0 go run cmd/main.go < input.txt
```

### Compilação nativa:

```Bash
go clean

go build -o capital-gains cmd/main.go

# opcional
export PROFIT_PERCENTAGE=0.30 \
export MAX_SELL_OPERATION_VALUE=25000.0

./capital-gains
```

Ele identifica automaticamente o sistema operacional e arquitetura de processador que ele está compilando, mas podemos setar manualmente:

```Bash
GOOS=darwin GOARCH=arm64 go build -o capital-gains main.go  # macos arm64
GOOS=windows GOARCH=amd64 go build -o capital-gains main.go # windows amd64
GOOS=linux GOARCH=amd64 go build -o capital-gains main.go   # linux amd64

./capital-gains
```

### Via Docker

```Bash
docker build -t capital-gains -f Docker.dockerfile .

docker run -i capital-gains < input.txt

# input customizável
docker run -e PROFIT_PERCENTAGE=0.30 -e MAX_SELL_OPERATION_VALUE=25000.0 -i capital-gains < input.txt
```

Run the tests:
```Bash

`````

## Decisões Arquiteturais

### Por que Go:

Minha decisão por utilizar GO foi pela natureza do problema de ser executado via CLI, onde ele lida muito bem com bulds nativas em múltiplos sistemas e com Docker,
não sendo necessária a utilização de interpretadores/VM's.

Outra foi mais pessoal de utilizar uma linguagem que não é tanto a minha zona de conforto como o Java.

### Dificuldades:

 - Sem dúvidas por utilizar apenas bibliotecas da própria linguagem, lidar com decoding/encoding foi desafiador, mas interessante de fazer funcionar!
 - Como lidar com floats para valores financeiros em Go? Resolvi utilizar o float64 por ser uma aplicação não produtiva e evitar o uso por enquanto de libs externas.
 - Pensar em uma estrutura que suporte mais de uma ação. `Esse não foi um requisito, mas achei válido essa extensibilidade`
