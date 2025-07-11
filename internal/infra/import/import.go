package importador

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/xuri/excelize/v2"

	"enube/internal/domain/model"
	"enube/internal/infra/db"
)

func ImportFromExcel(filepath string) error {
	// Log de erros
	logFilePath := filepath + ".log"
	logFile, err := os.Create(logFilePath)
	if err != nil {
		log.Println("Erro ao criar log de erro:", err)
		return err
	}
	defer logFile.Close()

	errorLog := log.New(logFile, "ERRO: ", log.LstdFlags)

	f, err := excelize.OpenFile(filepath)
	log.Printf("Lendo arquivo: %s\n", filepath)
	if err != nil {
		return fmt.Errorf("erro ao abrir arquivo: %w", err)
	}
	log.Printf("Arquivo aberto com sucesso: %s\n", filepath)
	log.Printf("Planilhas disponíveis: %v\n", f.GetSheetList())
	log.Printf("Lendo linhas da planilha 'Planilha1'...\n")
	rows, err := f.GetRows("Planilha1")
	if err != nil {
		return fmt.Errorf("erro ao ler planilha: %w", err)
	}
	log.Printf("Total de linhas lidas: %d\n", len(rows))
	log.Printf("Iniciando importação de dados...\n")

	// Caches
	partners := map[string]model.Partner{}
	customers := map[string]model.Customer{}
	products := map[string]model.Product{}
	subscriptions := map[string]model.Subscription{}
	meters := map[string]model.Meter{}
	var billingItems []model.BillingItem

	successCount := 0
	errorCount := 0

	for i, row := range rows {
		if i == 0 {
			continue
		}

		defer func() {
			if r := recover(); r != nil {
				errorCount++
				errorLog.Printf("linha %d: panic recover: %v\n", i, r)
			}
		}()

		// Caching: Partner
		partnerID := row[0]
		if _, ok := partners[partnerID]; !ok {
			partner := model.Partner{PartnerID: partnerID, Name: row[1]}
			if err := db.DB.FirstOrCreate(&partner, model.Partner{PartnerID: partnerID}).Error; err != nil {
				errorCount++
				errorLog.Printf("linha %d: partner: %v\n", i, err)
				continue
			}
			partners[partnerID] = partner
		}

		// Caching: Customer
		customerID := row[2]
		if _, ok := customers[customerID]; !ok {
			customer := model.Customer{
				CustomerID: customerID,
				Name:       row[3],
				DomainName: row[4],
				Country:    row[5],
			}
			if err := db.DB.FirstOrCreate(&customer, model.Customer{CustomerID: customerID}).Error; err != nil {
				errorCount++
				errorLog.Printf("linha %d: customer: %v\n", i, err)
				continue
			}
			customers[customerID] = customer
		}

		// Caching: Product
		productKey := row[9] + row[10]
		if _, ok := products[productKey]; !ok {
			product := model.Product{
				ProductID:      row[9],
				SkuID:          row[10],
				AvailabilityID: row[11],
				SkuName:        row[12],
				ProductName:    row[13],
				PublisherName:  row[14],
				PublisherID:    row[15],
			}
			if err := db.DB.FirstOrCreate(&product, model.Product{ProductID: row[9], SkuID: row[10]}).Error; err != nil {
				errorCount++
				errorLog.Printf("linha %d: product: %v\n", i, err)
				continue
			}
			products[productKey] = product
		}

		// Caching: Subscription
		subID := row[17]
		if _, ok := subscriptions[subID]; !ok {
			sub := model.Subscription{
				SubscriptionID: subID,
				Description:    row[16],
			}
			if err := db.DB.FirstOrCreate(&sub, model.Subscription{SubscriptionID: subID}).Error; err != nil {
				errorCount++
				errorLog.Printf("linha %d: subscription: %v\n", i, err)
				continue
			}
			subscriptions[subID] = sub
		}

		// Caching: Meter
		meterID := row[23]
		if _, ok := meters[meterID]; !ok {
			meter := model.Meter{
				MeterID:     meterID,
				Type:        row[21],
				Category:    row[22],
				SubCategory: row[24],
				Name:        row[25],
				Region:      row[26],
			}
			if err := db.DB.FirstOrCreate(&meter, model.Meter{MeterID: meterID}).Error; err != nil {
				errorCount++
				errorLog.Printf("linha %d: meter: %v\n", i, err)
				continue
			}
			meters[meterID] = meter
		}

		// Parsing
		chargeStart, _ := time.Parse("02-01-06", row[18])
		chargeEnd, _ := time.Parse("02-01-06", row[19])
		rateDate, _ := time.Parse("02-01-06", row[46])
		unitPrice, _ := strconv.ParseFloat(row[33], 64)
		quantity, _ := strconv.ParseFloat(row[34], 64)
		billTotal, _ := strconv.ParseFloat(row[36], 64)
		priceTotal, _ := strconv.ParseFloat(row[38], 64)
		effective, _ := strconv.ParseFloat(row[44], 64)
		rate, _ := strconv.ParseFloat(row[45], 64)
		creditPerc, _ := strconv.Atoi(row[49])
		creditPartnerPerc, _ := strconv.Atoi(row[50])

		billingItems = append(billingItems, model.BillingItem{
			InvoiceNumber:                 row[8],
			PartnerID:                     partners[partnerID].ID,
			CustomerID:                    customers[customerID].ID,
			ProductID:                     products[productKey].ID,
			SubscriptionID:                subscriptions[subID].ID,
			MeterID:                       meters[meterID].ID,
			ChargeStartDate:               chargeStart,
			ChargeEndDate:                 chargeEnd,
			UsageDate:                     row[20],
			ResourceLocation:              row[28],
			ResourceGroup:                 row[30],
			ResourceURI:                   row[31],
			ConsumedService:               row[29],
			ChargeType:                    row[32],
			Unit:                          row[27],
			UnitType:                      row[35],
			UnitPrice:                     unitPrice,
			Quantity:                      quantity,
			BillingPreTaxTotal:            billTotal,
			BillingCurrency:               row[37],
			PricingPreTaxTotal:            priceTotal,
			PricingCurrency:               row[39],
			EffectiveUnitPrice:            effective,
			PCToBCExchangeRate:            rate,
			PCToBCExchangeRateDate:        rateDate,
			EntitlementID:                 row[47],
			EntitlementDescription:        row[48],
			PartnerEarnedCreditPercentage: creditPartnerPerc,
			CreditPercentage:              creditPerc,
			CreditType:                    row[51],
			Tags:                          row[42],
			AdditionalInfo:                row[43],
			ServiceInfo1:                  row[40],
			ServiceInfo2:                  row[41],
		})

		successCount++

		if successCount%2000 == 0 {
			fmt.Printf("Importados: %d registros com sucesso...\n", successCount)
			if err := db.DB.CreateInBatches(billingItems, 200).Error; err != nil {
				errorLog.Printf("Erro ao salvar lote: %v\n", err)
			}
			billingItems = []model.BillingItem{}
		}
	}

	if len(billingItems) > 0 {
		if err := db.DB.CreateInBatches(billingItems, 100).Error; err != nil {
			errorLog.Printf("Erro ao salvar lote final: %v\n", err)
		}
	}

	log.Printf("Importação concluída: %d inserções com sucesso, %d erros.\n", successCount, errorCount)
	fmt.Printf("Arquivo de log salvo em: %s\n", logFilePath)

	return nil
}
