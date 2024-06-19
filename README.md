# CLOUDCOST
Ferramenta de extração e análise dos dados de billing em arquivos (CSV) dos provedores de nuvem como AWS, Azure, Google e etc.

# Overview

[ESTUDO DE CASO] Este utilitário realiza o processamento sobre os dados gerados em relatórios de billing de Cloud, permitindo a rápida extração de informações de relatórios CSV complexos como Cost Usage Report (CUR), Cost Explorer e outros.

O objetivo desta ferramenta é permitir a análise dos dados gerados no billing dos provedores de nuvem, usando concorrência/paralelismo para permitir a melhor eficiência na leitura de dados extensos e complexos. 

Cenário: O AWS Cost Usage Report ao longo de 30 dias pode gerar dezenas de arquivos CSV com dados de consumo de nuvem detalhados, gerando uma massa de dados que pode ultrapassar GB de armazenamento, para processar toda essa massa de dados e extrair informações rapidamente é necessário o uso de ferramentas de BI, muitas vezes inacessíveis ou de difícil uso. Neste tipo de cenário este utilitário permite a análise de toda essa massa de dados em segundos, facilitando a rápida extração de informações importantes como custo total, custo por produto consumido, custo individual por conta, resource ID, usage type, PTAX, fator de cobrança quando existente e outras informações.


# Uso

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
# Importando OVF ESXi

Host ESXi/vCenter > Deploy OVF Template > next... next... finish

# Download AWS ADS

Link download AWS ADS <a href="https://s3.us-west-2.amazonaws.com/aws.agentless.discovery.collector.bundle/releases/latest/ApplicationDiscoveryServiceAgentlessCollector.ova" target="_blank">ApplicationDiscoveryServiceAgentlessCollector.ova</a>


# Contato

Alex Manfrin <alexmanfrin@gmail.com>

Linkedin - https://www.linkedin.com/in/alexandermanfrin/



