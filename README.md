# Cloudcost
[![Go Reference](https://pkg.go.dev/badge/mitvix/cloudcost.svg)](https://pkg.go.dev/mitvix/cloudcost)

Ferramenta de extração e análise dos dados de billing em arquivos (CSV) dos provedores de nuvem como AWS, Azure, Google e etc.

## Overview

Este utilitário realiza o processamento sobre os dados gerados em relatórios de billing de Cloud, permitindo a rápida extração de informações de arquivos CSV complexos gerados a partir do Cost Usage Report (CUR) da AWS, Cost Explorer e outros.

O objetivo desta ferramenta é permitir a análise dos dados de billing dos provedores de nuvem, usando concorrência/paralelismo para permitir a melhor eficiência na leitura de dados extensos e complexos sem a necessidade de soluções de BI.

<details>

<summary>Cenário de uso</summary>

## Cenário

> O AWS Cost Usage Report pode gerar dezenas de arquivos CSV detalhados do consumo de nuvem, podendo gerar (GigaBytes) de dados em arquivos de texto. Para processar essa massa de dados e extrair informações rapidamente é necessário o uso de ferramentas de BI, muitas vezes inacessíveis ou de difícil uso. Neste cenário este utilitário permite a análise de toda a massa de dados em segundos, facilitando a rápida extração das informações mais importantes como custos dos produtos, uso de recursos, total por conta, resource ID, usage type, PTAX, fator de cobrança quando existente, além de outras customizações e informações importantes.


</details>

## [ESTUDO DE CASO]
Software criado _(nas raras horas vagas)_ para estudo e análise da línguagem Go (Golang) disponível em [go.dev](https://go.dev). Línguagem de programação opensource criada por [Rob Pike](https://pt.wikipedia.org/wiki/Rob_Pike), [Robert Griesemer](https://en.wikipedia.org/wiki/Robert_Griesemer) e [Ken Thompson](https://pt.wikipedia.org/wiki/Ken_Thompson) nos laboratórios do Google em meados de 2007 e liberado sob licença opensource BSD em 2009.

<details>

<summary>Sobre Go e por que essa linguagem</summary>

### Golang

Go foi projetado inicialmente com o objetivo de substituir projetos em C e C++ dentro do Google, por isso possui características simílares a essas línguagens, incluindo sua síntaxe, mas com abstrações voltadas a simplicidade e legibilidade, além de uma forte combinação de suporte a concorrência e desempenho. Sua estrutura automática de gerenciamento de memória (Garbage Collector) facilita a vida do desenvolvedor, mas gera overhead que a deixa pouco atrás em performance quando comparada a C, C++ e Rust, porém, muito a frente em desempenho em relação a Python, Java, PHP e etc. 

E mesmo perdendo em performance para Rust e C++, Go se tornou uma línguagem equilibrada que combina estruturas de baixo nível de C com a usabilidade do mundo moderno e sem o pesadelo da Orientação a Objetos, fazendo dela uma línguagem de programação simples, completa e perfeita para o uso em APIs, Micro serviços, Web Dev, Cloud e etc. 

Dentre os principais projetos escritos em Go, temos Kubernetes, kubectl, Minikube, Docker e outros. Veja mais em [https://go.dev/solutions/cloud#use-case](https://go.dev/solutions/cloud#use-case)

### Por que Go ?

_Seria em Rust, mas a partir do capítulo 5 Rust se tornou complexo demais_

O objetivo foi de focar esforços no aprendizado de uma línguagem moderna para pequenos projetos paralelos, em Go foi possível encontrar uma curva de aprendizado curta, pois sua semântica segue padrões já conhecidos de outras línguagens como C, mas sem as dores do uso de uma línguagem nosafe. Aliado a capacidade de criar códigos usando concorrência e paralelismo de forma rápida e simples, Go demonstrou ser a opção mais adequada para o objetivo de percorrer gigabytes de arquivos em busca de padrões para extrair informações de forma rápida e precisa sem a necessidade de um doutorado para isso.

</details>

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
| Data de início |  Data final |  Cotação do Dólar e data | Fator de cobrança |  Total em dólar |  Total em reais | 

Exemplo
| INÍCIO     |  FIM        |  PTAX              | FATOR  |  CONSUMO USD |  CONSUMO BRL | 
| --- | --- | --- | --- | --- | --- |
| 01-03-2024 | 31-03-2024 | 4.9962 2024-03-28 |  1 | US$ 4094.42 | R$ 22911.31 |

Outras saídas

* Resource Id
* Usage Type
* Market Place
* Savings
* Credits
* Support

### Exemplos

Retorna o UsageType e informações sobre ExtendedSupport e RDS
```
cloudcost path_file_csv/ --usagetype --search ExtendedSupport,RDS
```

Retorna o Resource ID e todos os recursos que contém snapshots
```
cloudcost path_file_csv/ --resourceid --search snapshot
```
Retorna somente os dados das contas X e Y
```
cloudcost report_file.csv --account 123456789090,098765432109
```

Tela com saída padrão de análise do Cost Usage Report

![Tela de saída da análise de arquivos do AWS CUR](https://github.com/mitvix/cloudcost/assets/12394000/c3365a29-b794-4d9b-9f7d-c676c9d4ec68)

Tela com outros tipos de saídas

![screen_savings](https://github.com/mitvix/cloudcost/assets/12394000/715cc78e-56bb-4e7d-8616-82a58d0a812a)

![screen_resourceid](https://github.com/mitvix/cloudcost/assets/12394000/2f033df8-78c4-43aa-a775-349644c110fb)

![screen_usagetype](https://github.com/mitvix/cloudcost/assets/12394000/14227dc9-593c-4dbf-8cca-0e9134b7e0e8)


## Benchmark

Tempo médio de processamento de arquivos csv

* 1.7GB (37 arquivos * 52MB) em 5.8s (1.960.024 linhas)
* 757MB (28 arquivos * 32MB) em 3.28s (1.211.316 linhas)
* 11GB  (1 arquivo * 11GB)   em 1m14s (13.438.628 linhas) w/ GC control
* 424MB (1 arquivo * 424MB)  em 2.14s (427.229 linhas)
* 115MB (1 arquivo * 115MB)  em 0.56s (144.117 linhas)

> Hardware 12 × Intel® Core™ i7-9750H CPU @ 2.60GH, 32GB Kernel 6.6.32 (64-bit);
> GC Control: Garbage Collector limitado até 12GB default (--memlimit para alterar)

## Pacotes utilizados, consultas API e outros

- Uso de pacotes nativos da línguagem somente
- RFC4180 CSV https://www.ietf.org/rfc/rfc4180.txt
- API BC (PTAX) https://olinda.bcb.gov.br/olinda/servico/PTAX/versao/v1/aplicacao#!/recursos
- ASCII documentation https://gist.github.com/fnky/458719343aabd01cfb17a3a4f7296797
- Debug Memory https://golang.org/pkg/runtime/#MemStats  

## GOAL
Lista de objetivos para o programa

```
Carregamento de gigabytes em segundos            :  97%
Análise de custos por conta                      :  70%
Manutenção externa de padrões de cabeçalho csv   :  1%
Multiplicação de fator / ptax por produto        :  1%
Uso de fator de consumo/marketplace customizado  :  55%
Uso de ptax customizada                          :  15%
Análise de relatórios AWS CUR                    :  71%
Análise de relatórios terceiros                  :  81%
Análise de relatórios Cost Explorer              :  49% w/ issues
Análise de relatórios Microsoft                  :  1%
Análise de relatórios Google                     :  1%
Funções de pesquisas customizadas --search       :  63%
Função de coleta de dólar ptax                   :  77% w/ issues
Função export para TXT                           :  4%
Função export para PDF                           :  0%

```

## Nota do Autor e Contato

* _IMPORTANTE O autor não é programador e não deseja ser um, programar é uma arte, um hobby como dar um Dolio Tchagui ou ir para cozinha para fazer uma Paeja ou um Moti! Então, ignore código fora das melhores práticas e típicas de um newbie, por isso é GPL, encontrou um erro, arrume! provavelmente há vários deles..._

Alex Manfrin <mitvix@hotmail.com>
Linkedin - https://www.linkedin.com/in/alexandermanfrin/

