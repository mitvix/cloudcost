# CLOUDCOST
Ferramenta de extração e análise dos dados de billing em arquivos (CSV) dos provedores de nuvem como AWS, Azure, Google e etc.

## Overview

Este utilitário realiza o processamento sobre os dados gerados em relatórios de billing de Cloud, permitindo a rápida extração de informações de arquivos CSV complexos gerados a partir do Cost Usage Report (CUR) da AWS, Cost Explorer e outros.

O objetivo desta ferramenta é permitir a análise dos dados de billing dos provedores de nuvem, usando concorrência/paralelismo para permitir a melhor eficiência na leitura de dados extensos e complexos sem a necessidade de soluções de BI.

**Cenário**: O AWS Cost Usage Report pode gerar dezenas de arquivos CSV detalhados do consumo de nuvem, podendo gerar (GigaBytes) de dados em arquivos. Para processar essa massa de dados e extrair informações rapidamente é necessário o uso de ferramentas de BI, muitas vezes inacessíveis ou de difícil uso. Neste cenário este utilitário permite a análise de toda a massa de dados em segundos, facilitando a rápida extração das informações mais importantes como custos dos produtos, uso de recursos, total por conta, resource ID, usage type, PTAX, fator de cobrança quando existente e outras informações.

## [ESTUDO DE CASO]
Software criado _(nas raras horas vagas)_ para estudo e análise da línguagem Go (Golang) disponível em [go.dev](https://go.dev). Línguagem de programação opensource criada por [Rob Pike](https://pt.wikipedia.org/wiki/Rob_Pike), [Robert Griesemer](https://en.wikipedia.org/wiki/Robert_Griesemer) e [Ken Thompson](https://pt.wikipedia.org/wiki/Ken_Thompson) nos laboratórios do Google em meados de 2007 e liberado sob licença opensource BSD em 2009.

Go foi projetado inicialmente com o objetivo de substituir projetos em C e C++ dentro do Google, por isso possui características simílares a essas línguagens, incluindo sua síntaxe, mas com abstrações voltadas a simplicidade e legibilidade, além de uma forte combinação de suporte a concorrência e desempenho. Sua estrutura automática de gerenciamento de memória (Garbage Collector) facilita a vida do desenvolvedor, mas gera overhead que a deixa pouco atrás em performance quando comparada a C, C++ e Rust, porém, muito a frente em desempenho em relação a Python, Java, PHP e etc. E mesmo perdendo em performance para Rust e C++, Go se tornou uma línguagem equilibrada que combina estruturas de baixo nível de C com a usabilidade do mundo moderno e sem o pesadelo da Orientação a Objetos, fazendo dela uma línguagem de programação simples, completa e perfeita para o uso em APIs, Micro serviços, Web Dev, Cloud e etc. Dentre os principais projetos escritos em Go, temos Kubernetes, kubectl, Minikube, Docker e outros. Veja mais em [https://go.dev/solutions/cloud#use-case](https://go.dev/solutions/cloud#use-case)

## Instalação

Para instalação do Go siga o passo a passo disponível em [https://go.dev/doc/tutorial/getting-started#install](go.dev).

Instalação
```
git clone https://github.com/mitvix/cloudcost
cd cloudcost
go build -o cloudcost main.go
```

Opcional (Compilação para Windows)
```
export GOOS=windows
go build -o cloudcost.exe main.go
```

Opcional (copy to bin path)
```
sudo cp cloudcost /usr/local/bin
```

## Uso

Mostrar opções disponíveis do programa
```
cloudcost --help
```

Argumentos:

* -path, --path `Diretório com arquivos CSV ou arquivo .csv`
* -header, --header `Mostra o cabeçalho do arquivo CSV (requer --path)`
* -account, --account `Filtro de análise por conta ex: --account 868884350453,443786768377 (requer --path)`
* -fee, --fee `Define fator de consumo padrão ex: --fee 1.09 (requer --path)`
* -feemp, --feemp `Define fator de consumo Market Place ex: --feemp 1.7550 (requer --path)`
* -marketplace, --marketplace `Mostra os detalhes de recursos do Market Place`
* -memlimit, --memlimit `Define max memory MB em uso - tenta controlar GC e pode gerar lentidão`
* -resourceid, --resourceid `Mostra detalhes de recursos por ID/arn (ResourceID)`
* -search, --search `Faz busca nos relatórios min. 2 caracteres`
* -usagetype, --usagetype `Mostra detalhes do tipo do recurso (UsageType)`
* -version, --version `Mostra informações sobre a versão e sai`

## Saídas

Consumo

| INÍCIO     |  FIM        |  PTAX              | FATOR  |  CONSUMO USD |  CONSUMO BRL | 
| --- | --- | --- | --- | --- | --- |
| Data de início |  Data final |  Cotação do Dólar e data | Fator de cobrança |  Total em dólar |  Todal em reais | 

Exemplo
| INÍCIO     |  FIM        |  PTAX              | FATOR  |  CONSUMO USD |  CONSUMO BRL | 
| --- | --- | --- | --- | --- | --- |
| 01-03-2024 | 31-03-2024 | 4.9962 2024-03-28 |  1 | US$ 4094.42 | R$ 22911.31 |

Outras saídas

* Resource Id
* Usage Type
* Market Place
* Credits
* Support

# Nota do Autor e Contato

* _IMPORTANTE O autor não é programador e não deseja ser um, programar é uma arte, um hobby como dar um Dolio Tchagui ou ir para cozinha para fazer uma Paeja ou um Moti! Então, ignore código fora das melhores práticas e típicas de um newbie, por isso é GPL, encontrou um erro, arrume!_

Alex Manfrin <mitvix@hotmail.com>
Linkedin - https://www.linkedin.com/in/alexandermanfrin/

