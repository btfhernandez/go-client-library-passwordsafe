package main

import (
	"fmt"
	"time"

	"github.com/bthernandez/go-client-library-passwordsafe/api/authentication"
	"github.com/bthernandez/go-client-library-passwordsafe/api/entities"
	logging "github.com/bthernandez/go-client-library-passwordsafe/api/logging"
	managed_accounts "github.com/bthernandez/go-client-library-passwordsafe/api/managed_account"
	"github.com/bthernandez/go-client-library-passwordsafe/api/utils"

	//"os"

	backoff "github.com/cenkalti/backoff/v4"
	"go.uber.org/zap"
)

// main funtion
func main() {

	// create a zap logger
	//logger, _ := zap.NewProduction()
	logger, _ := zap.NewDevelopment()

	// create a zap logger wrapper
	zapLogger := logging.NewZapLogger(logger)

	apiUrl := "https://jury2310.ps-dev.beyondtrustcloud.com:443/BeyondTrust/api/public/v3"
	clientId := "6138d050-e266-4b05-9ced-35e7dd5093ae"
	clientSecret := "71svdPLh2AR97sPs5gfPjGjpqSUxZTKSPmEvvbMx89o="
	separator := "/"
	certificate := ""
	certificateKey := ""
	clientTimeOutInSeconds := 30
	verifyCa := true
	retryMaxElapsedTimeMinutes := 2
	maxFileSecretSizeBytes := 5000000

	backoffDefinition := backoff.NewExponentialBackOff()
	backoffDefinition.InitialInterval = 1 * time.Second
	backoffDefinition.MaxElapsedTime = time.Duration(retryMaxElapsedTimeMinutes) * time.Second
	backoffDefinition.RandomizationFactor = 0.5

	//certificate = os.Getenv("CERTIFICATE")
	//certificateKey = os.Getenv("CERTIFICATE_KEY")

	// Create an instance of ValidationParams
	params := utils.ValidationParams{
		ClientID:                   clientId,
		ClientSecret:               clientSecret,
		ApiUrl:                     &apiUrl,
		ClientTimeOutInSeconds:     clientTimeOutInSeconds,
		Separator:                  &separator,
		VerifyCa:                   verifyCa,
		Logger:                     zapLogger,
		Certificate:                certificate,
		CertificateKey:             certificateKey,
		RetryMaxElapsedTimeMinutes: &retryMaxElapsedTimeMinutes,
		MaxFileSecretSizeBytes:     &maxFileSecretSizeBytes,
	}

	// validate inputs
	errorsInInputs := utils.ValidateInputs(params)

	if errorsInInputs != nil {
		return
	}

	// creating a http client
	httpClientObj, _ := utils.GetHttpClient(clientTimeOutInSeconds, verifyCa, certificate, certificateKey, zapLogger)

	// instantiating authenticate obj, injecting httpClient object
	authenticate, _ := authentication.Authenticate(*httpClientObj, backoffDefinition, apiUrl, clientId, clientSecret, zapLogger, retryMaxElapsedTimeMinutes)

	// authenticating
	_, err := authenticate.GetPasswordSafeAuthentication()
	if err != nil {
		return
	}

	// instantiating managed account obj
	manageAccountObj, _ := managed_accounts.NewManagedAccountObj(*authenticate, zapLogger)

	account := entities.AccountDetails{
		AccountName:                       "Test2014R",
		Password:                          "Hol",
		DomainName:                        "exampleDomain",
		UserPrincipalName:                 "user@example.com",
		SAMAccountName:                    "samAccount",
		DistinguishedName:                 "CN=example,CN=Users,DC=domain,DC=com",
		PrivateKey:                        "privateKey",
		Passphrase:                        "passphrase",
		PasswordFallbackFlag:              true,
		LoginAccountFlag:                  false,
		Description:                       "Sample account for testing",
		ApiEnabled:                        true,
		ReleaseNotificationEmail:          "notify@example.com",
		ChangeServicesFlag:                false,
		RestartServicesFlag:               false,
		ChangeTasksFlag:                   true,
		ReleaseDuration:                   300000,
		MaxReleaseDuration:                300000,
		ISAReleaseDuration:                180,
		MaxConcurrentRequests:             5,
		AutoManagementFlag:                false,
		DSSAutoManagementFlag:             false,
		CheckPasswordFlag:                 true,
		ResetPasswordOnMismatchFlag:       false,
		ChangePasswordAfterAnyReleaseFlag: true,
		ChangeFrequencyType:               "first",
		ChangeFrequencyDays:               1,
		ChangeTime:                        "14:00",
		NextChangeDate:                    "2023-12-01", // Cambiar a `time.Time` si es necesario
		UseOwnCredentials:                 true,
		ChangeWindowsAutoLogonFlag:        true,
		ChangeComPlusFlag:                 false,
		ObjectID:                          "uniqueObjectID",
	}

	// creating a managed account in system_integration_test managed system.
	createResponse, err := manageAccountObj.ManageAccountCreateFlow("system_integration_test", account)

	if err != nil {
		zapLogger.Debug(fmt.Sprintf(" %v", err))
		return
	}

	zapLogger.Debug(fmt.Sprintf("Created Managed Account: %v", createResponse.AccountName))

	// signing out
	_ = authenticate.SignOut()

}
