package global

import "time"

/*
 * Main configuration
 * Constants, Strings and global variables
 */

// Max limits
const Max_limitfiles = 512 // limit files to analysis 0 = ilimited caution!
const Max_memlimit = 12e+9 // 12e+9 Example 12GB (8GB = 8e+9) 0 == ilimited
const Max_charprint = 100  // max characters to print product name

var Memory_alert int64 = 12000         // Memory alert > (12GB = 12e+9) = red(MemSys)
const Ptaxlayout string = "01-02-2006" // 01 (month) 02 (day) 06 (year)
const Datelayout string = "02-01-2006" // BC API layout 02 (day) 01 (month) 06 (year)
const PtaxTimeout time.Duration = 3    // http.Client {Timeout: T * time.Second }

// Default fee value
const DefaultFee float64 = 1 // Default fee in non BRL reports
const PipeLine int = 30      // Default pipe line to show
var ShowPipe bool = true     // Default configuration for pipe line = true
var UseColor bool = true     // UseColors determines if package colors will emit colors.

// ProgressBar
const Barenable bool = true     // Enable/disable Progress Bar (true|false)
const Barsize int = 39          // ShowProgressBar size
const Barchar string = "\u2587" // ShowProgressBar character \u220E
const BarcharWin string = "#"
const InnerBar string = "[ ....................................... ]" // 39 positions!

// Main Strings and Values
const (
	Version  string = "0.0.23"
	Program  string = "Cloud Cost Report Reader"
	Codename string = "codename Stallman"
	Aws_cexp string = "AWS Cost Explorer"
	Aws_cur  string = "AWS Cost Usage Report"
	Cmp      string = "Cmp"
	Azure    string = "Microsoft Azure"
	Google   string = "Google GCP"

	// define report cloud
	RepAws    string = "aws"
	RepOci    string = "oci"
	RepAzure  string = "azure"
	RepGoogle string = "google"
	RepHuawei string = "huawei"

	Src_cexp string = "LinkedAccountId" // aws cost explorer US
	Src_cost string = "InvoiceID"
	Src_Cmp  string = "Usage Account"
	Src_cur  string = "lineItem/UsageAccountId"
	Src_azur string = "SubscriptionId"
	Src_gcp  string = "billing_account_id"

	Progdesc    string = "Busca padrões em arquivos de relatórios Multicloud e filtra dados de billing"
	License     string = "\nLicença GPLv3+: GNU GPL versão 3 ou superior <https://gnu.org/licenses/gpl.html>.\nEste é um software livre: você é livre para alterá-lo e redistribuí-lo.\nNÃO HÁ GARANTIAS, na máxima extensão permitida por lei. \n\nEscrito por Alexander Manfrin <mitvix@hotmail.com> em"
	Flagdirpath string = "path"
	Flagaccount string = "account"
	Flagfee     string = "fee"
	Flagfeemp   string = "feemp"
	Flagptax    string = "ptax"
	Flagptaxmp  string = "ptaxmp"
	Flagsheader string = "header"
	Flagusgtype string = "usagetype"
	Flagresrcid string = "resourceid"
	Flagmarkplc string = "marketplace"
	Flagmemlimt string = "memlimit"
	Flagresrcgr string = "resourcegroup"
	Flagnopipe  string = "nopipe"
	Flagsearch  string = "search"
	Flagversion string = "version"

	Msg_flagdir string = "Diretório com arquivos CSV ou arquivo .csv"
	Msg_flagacc string = "Filtro de análise por conta ex: --account 868884350453,443786768377 (requer --path)"
	Msg_flagfee string = "Define fator de consumo padrão ex: --fee 1.09 (requer --path)"
	Msg_flfeemp string = "Define fator de consumo Market Place ex: --feemp 1.7550 (requer --path)"
	Msg_flagptx string = "Define PTAX de consumo padrão ex: --ptax 4.9962 (requer --path)"
	Msg_flptxmp string = "Define PTAX de consumo Market Place ex: --ptaxmp 5.19 (requer --path)"

	Msg_flagptax string = "Define PTAX padrão manualmente (requer --path)"
	Msg_flaghed  string = "Mostra o cabeçalho do arquivo CSV (requer --path)"
	Msg_flagrsc  string = "Mostra detalhes do tipo do recurso (UsageType)"
	Msg_flagrsg  string = "Mostra custos por Resource Group em relatórios Microsoft Azure"
	Msg_flagrid  string = "Mostra detalhes de recursos por ID/arn (ResourceID)"
	Msg_flagmkl  string = "Mostra os detalhes de recursos do Market Place"
	Msg_memlimt  string = "Define max memory MB em uso - tenta controlar GC e pode gerar lentidão"
	Msg_nopipe   string = "Desativa a pausa na visualização e formatações de texto"
	Msg_flagsch  string = "Faz busca nos relatórios min. 2 caracteres (requer --path)"
	Msg_version  string = "Mostra informações sobre a versão e sai"

	Msg_contrt     string = "Contrato:"
	Msg_rangstr    string = "Início"
	Msg_rangend    string = "Fim:"
	Msg_accout     string = "Accounts:"
	Msg_produt     string = "Produtos:"
	Msg_source     string = "Billing:"
	Msg_rsrctype   string = "Usage type:"
	Msg_resrcid    string = "Resource IDs:"
	Msg_resrgrp    string = "Resource Group:"
	Msg_total      string = "Total:"
	Msg_usage      string = "Consumo:"
	Msg_ptaxglb    string = "Ptax:"
	Msg_mktpl      string = "Market Place:"
	Msg_feeglb     string = "Fator Consumo:"
	Msg_credit     string = "Créditos:"
	Msg_support    string = "Suporte:"
	Msg_savings    string = "Savings:"
	Msg_feerr      string = "-"
	Msg_valid      string = "Validação do cálculo"
	Msg_ok         string = "[OK]"
	Msg_error      string = "ERRO!"
	Msg_vdiff      string = "Valores diferentes!"
	Msg_fileopen   string = "Erro ao abrir o arquivo"
	Msg_fileread   string = "Erro ao ler linha do arquivo "
	Msg_fileerro   string = "Arquivo ou diretório não encontrado"
	Msg_fileform   string = "Arquivo mal formatado"
	Msg_prsenter   string = "Pressione ENTER para continuar"
	Msg_anotfoud   string = "Conta não encontrada"
	Msg_snotfoud   string = "Resultado não encontrado"
	Msg_connon     string = "Consulta Ptax Online \u2713"
	Msg_connoff    string = "Timeout ou sem conexão para coletar Ptax \u274C"
	Msg_file       string = "Arquivo"
	Msg_sumtotal   string = "Soma Total"
	Msg_sumprodt   string = "Soma por produto"
	Msg_extcsv     string = ".csv"
	Msg_maxreached string = "Limite máximo de arquivos para análise alcançado"

	Msg_SymblBR string = "R$"
	Msg_SymblUS string = "US$"
	SymbolBRL   string = "BRL"
	SymbolUSD   string = "USD"
)

var (
	Msg_exemp  string   = "Exemplos:"
	Msg_endex  string   = "Para mais exemplos, acesse a página do projeto(1):"
	Msg_github string   = "(1) https://github.com/mitvix/cloudcost"
	Examples   []string = []string{
		"--path=/reports/csv/",
		"--account=621882272082 --path=/reports/csv/",
		"--account=621882272082,210919792179 --path=/reports/csv/",
		"/relatorios/csv --search db.m5.large,Lambda,WorkSpaces",
		"--path relatorios/csv/file.csv --fee 1.09 --feemp 1.7500",
		"--path relatorios/csv/file.csv --usagetype",
		"--path relatorios/csv/file.csv --resourceid",
	}
)

var (
	Msg_withfee string = "Conversão dolar s/ fator"
	Msg_without string = "Fator adicionado"
)

// Field prefix to search for usages  (BE CAREFUL!!!)
var (
	Marketplace []string = []string{"MP:", "SoftwareUsage", "Marketplace"} // strings to search marketplace
	Credits     []string = []string{"Credit"}                              // strings to search credits
	Support     []string = []string{"AWS Support (Enterprise)"}            // strings to search support
	Savings     []string = []string{"Savings"}                             // strings to search Savings

)

// Field Position on CSV files (Cmp, Cmp GOV, CUR, Cost Explorer)
var (
	StartDate     = []string{"Start Date", "bill/BillingPeriodStartDate", "BillingPeriodStartDate"}
	EndDate       = []string{"End Date", "bill/BillingPeriodEndDate", "BillingPeriodEndDate"}
	CompanyName   = []string{"Company Name", "bill/InvoicingEntity", "PayerAccountName"}
	UsageAccount  = []string{"Usage Account", "lineItem/UsageAccountId", "LinkedAccountName"}
	ProductName   = []string{"Product Name", "product/ProductName", "ProductName"}
	UsageType     = []string{"Usage Type", "UsageType", "lineItem/UsageType"} // 1:1 "Cmp","CExplorer", CUR...  TO-DO list Azure, Google
	ResourceIdent = []string{"Resource Identifier", "lineItem/ResourceId"}
	ResourceCost  = []string{"Resource Cost"}
	FinalCost     = []string{"Final Cost", "Final Price (R$)", "CostBeforeTax", "lineItem/UnblendedCost"}
	CurrencyCode  = []string{"lineItem/CurrencyCode", "CurrencyCode"}
	ReportCloud   = []string{"aws", "azure", "google", "huawei", "oci"}
)

var (
	// PTAX main variables
	PtaxAPIUrl    string = "https://olinda.bcb.gov.br/olinda/servico/PTAX/versao/v1/odata/CotacaoMoedaPeriodo(moeda=@moeda,dataInicial=@dataInicial,dataFinalCotacao=@dataFinalCotacao)?@moeda='%s'&@dataInicial='%s'&@dataFinalCotacao='%s'&$top=1&$orderby=dataHoraCotacao%%20desc&$format=json&$select=cotacaoVenda,dataHoraCotacao"
	PtaxRespError string = "PTAX: Erro ao ler a resposta da PTAX API:"
	PtaxDataError string = "Ptax Erro: Não foi possível determinar a data de fechamento"
	PtaxPtaxError string = "Ptax Erro: Não foi possível encontrar o dólar ptax"
	PtaxJsonError string = "PTAX: Erro unmarshalling JSON:"
	PtaxSyntError string = "syntax error at byte"
	PtaxIndexMapD string = "dataHoraCotacao"
	PtaxIndexMapV string = "cotacaoVenda"
)

var (
	MarketP_Header  string = "PTAX\t\tFATOR\t\tMP USD\t\tMP BRL"
	Account_Header  string = "\tCONTA\tUSD\t\tBRL\t\n"
	Product_Header  string = "\tRECURSO\tUSD\t\tBRL\t\n"
	Product_vHeader string = "\tUSD\t\tBRL\t\tRECURSO\n"
	Resource_Header string = "INÍCIO\t\tFIM\t\tPTAX\t\tFATOR\t\tCONSUMO USD\t\tCONSUMO BRL"
	Detail_Header   string = "TOTAL\tTIPO"
	Total_Header    string = "TOTAL USD\t\tTOTAL BRL"
)

// leave the last line empty
