// +build test

package service_test

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	. "github.com/fractalpal/eventflow-example/payment/adapters/service"
	"github.com/fractalpal/eventflow-example/payment/app"
	"testing"
)

type ServiceTestSuite struct {
	suite.Suite
	VariableThatShouldStartAtFive int
	IDProvider                    *MockedIDProvider
	Validator                     *MockedValidator
	Repository                    *MockedRepository
	Service                       app.Service
}

func (suite *ServiceTestSuite) SetupTest() {
	suite.IDProvider = new(MockedIDProvider)
	suite.Validator = new(MockedValidator)
	suite.Repository = new(MockedRepository)
	suite.Service = New(suite.IDProvider, suite.Repository, suite.Validator)
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestServiceTestSuite(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}

func (suite *ServiceTestSuite) Test_CreatePaymentSuccess() {
	var payment app.Payment
	id := "ee3e97ba-04bf-4422-a112-08c6ca89effb"
	ctx := context.TODO()
	suite.IDProvider.On("ID").Return(id)
	payment.ID = id
	suite.Validator.On("Validate", ctx, payment).Return(nil)
	suite.Repository.On("Save", ctx, payment).Return(nil)

	created, err := suite.Service.Create(ctx, payment)

	assert.NoError(suite.T(), err, "expected no error")
	assert.Equal(suite.T(), id, created.ID)
}

func (suite *ServiceTestSuite) Test_DeletePaymentSuccess() {
	id := "ee3e97ba-04bf-4422-a112-08c6ca89effb"
	ctx := context.TODO()
	suite.Repository.On("Delete", ctx, id).Return(nil)

	err := suite.Service.Delete(ctx, id)

	assert.NoError(suite.T(), err, "expected no error")
}

func (suite *ServiceTestSuite) Test_UpdateBeneficiarySuccess() {
	party := app.ThirdParty{}
	ctx := context.TODO()
	suite.Repository.On("UpdateBeneficiary", ctx, party).Return(nil)

	err := suite.Service.UpdateThirdParty(ctx, party, "beneficiary")

	assert.NoError(suite.T(), err, "error after service UpdateThirdParty")
}

func (suite *ServiceTestSuite) Test_UpdateBeneficiaryNoRecords() {
	party := app.ThirdParty{}
	ctx := context.TODO()
	suite.Repository.On("UpdateBeneficiary", ctx, party).Return(nil)

	err := suite.Service.UpdateThirdParty(ctx, party, "beneficiary")

	assert.NoError(suite.T(), err, "error after service UpdateThirdParty")
}

func (suite *ServiceTestSuite) Test_UpdateDebtorSuccess() {
	party := app.ThirdParty{}
	ctx := context.TODO()
	suite.Repository.On("UpdateDebtor", ctx, party).Return(nil)

	err := suite.Service.UpdateThirdParty(ctx, party, "debtor")

	assert.NoError(suite.T(), err, "error after service UpdateThirdParty")
}
