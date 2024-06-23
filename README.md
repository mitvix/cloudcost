# CLOUDCOST
Ferramenta de extração e análise dos dados de billing em arquivos (CSV) dos provedores de nuvem como AWS, Azure, Google e etc.

## Overview

Este utilitário realiza o processamento sobre os dados gerados em relatórios de billing de Cloud, permitindo a rápida extração de informações de arquivos CSV complexos gerados a partir do Cost Usage Report (CUR) da AWS, Cost Explorer e outros.

O objetivo desta ferramenta é permitir a análise dos dados de billing dos provedores de nuvem, usando concorrência/paralelismo para permitir a melhor eficiência na leitura de dados extensos e complexos sem a necessidade de soluções de BI.

**Cenário**: O AWS Cost Usage Report pode gerar dezenas de arquivos CSV detalhados do consumo de nuvem, podendo gerar (GigaBytes) de dados em arquivos. Para processar essa massa de dados e extrair informações rapidamente é necessário o uso de ferramentas de BI, muitas vezes inacessíveis ou de difícil uso. Neste cenário este utilitário permite a análise de toda a massa de dados em segundos, facilitando a rápida extração das informações mais importantes como custos dos produtos, uso de recursos, total por conta, resource ID, usage type, PTAX, fator de cobrança quando existente e outras informações.

## [ESTUDO DE CASO]
Software criado _(nas raras horas vagas)_ para estudo e análise da línguagem Go (Golang) disponível em [go.dev](https://go.dev). Línguagem de programação opensource criada por [Rob Pike](https://pt.wikipedia.org/wiki/Rob_Pike), [Robert Griesemer](https://en.wikipedia.org/wiki/Robert_Griesemer) e [Ken Thompson](https://pt.wikipedia.org/wiki/Ken_Thompson) nos laboratórios do Google em meados de 2007 e liberado sob licença opensource BSD em 2009.

Go foi projetado inicialmente com o objetivo de substituir projetos em C e C++ dentro do Google, por este motivo possui características muito simílares a essas línguagens, incluindo parte de sua síntaxe, mas com abstrações mais voltadas a simplicidade e legibilidade, além de uma forte combinação de suporte a concorrência e desempenho. Sua estrutura automática de gerenciamento de memória (Garbage Collector) facilita a vida do desenvolvedor, embora este fato propicie a perda de performance em certos casos, a deixando pouco atrás em performance quando comparada a C, C++ e Rust, porém, muito a frente em desempenho em relação a línguagens como Python e Java. E mesmo perdendo em performance para Rust e C++, Go se tornou uma línguagem equilibrada que combina as estruturas de baixo nível de C com a usabilidade do mundo moderno e cross platform, sem o pesadelo da Orientação a Objetos, fazendo dela uma línguagem de programação simples, completa e perfeita para o uso em APIs, Micro serviços, Web Development, Cloud, CLI e outros. Dentre os principais projetos escritos em Go, temos Kubernetes, kubectl, Minikube e etc. Veja mais em [https://go.dev/solutions/cloud#use-case](https://go.dev/solutions/cloud#use-case)

## Uso

Em sistemas Linux utilize a versão shell do script ou a versão compilada em go

Em sistemas Windows utilize a versão compilada em go

Antes de continuar baixe a última versão do VMware-ovftool-XXX-lin.x86_64.bundle.

  https://customerconnect.vmware.com/downloads/get-download?downloadGroup=OVFTOOL441

Baixe o script shell/convertads2ova.sh

[https://github.com/mitvix/convertads2ova.git](https://github.com/mitvix/convertads2ova/archive/refs/heads/main.zip)

Shell 

```
Use ./convertadsova.sh nome_do_virtual_appliance.ova
```

Go 

Em sistemas Linux use: 
```
./convertadsova nome_do_virtual_appliance.ova
```
Em sistemas Windows use: 
```
./convertadsova.exe nome_do_virtual_appliance.ova
```
A partir do código fonte convertads2ova.go use:
```
go run convertads2ova.go nome_do_virtual_appliance.ova
```

Exemplo shell:
```
  $ chmod +x convertadsova.sh

  $ ./convertadsova.sh ApplicationDiscoveryServiceAgentlessCollector.ova
```
## Importando OVF ESXi

Host ESXi/vCenter > Deploy OVF Template > next... next... finish

## Download AWS ADS

Link download AWS ADS <a href="https://s3.us-west-2.amazonaws.com/aws.agentless.discovery.collector.bundle/releases/latest/ApplicationDiscoveryServiceAgentlessCollector.ova" target="_blank">ApplicationDiscoveryServiceAgentlessCollector.ova</a>

## Contato e Nota do Autor

* _IMPORTANTE O autor não é programador e não deseja ser um, programar é uma arte, um hobby como ir para cozinha para fazer uma Paeja ou um Moti! Então, ignore código fora das melhores práticas e típicas de um newbie, por isso é GPL, encontrou um erro, arrume!_

Alex Manfrin <mitvix@hotmail.com>
Linkedin - https://www.linkedin.com/in/alexandermanfrin/

