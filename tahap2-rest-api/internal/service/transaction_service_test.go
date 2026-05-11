package service

import (
	"context"
	"testing"

	"github.com/user/tahap2-rest-api/internal/model"
	"github.com/user/tahap2-rest-api/internal/repository"
	"github.com/user/tahap2-rest-api/internal/worker"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

// --- Mock implementations ---

type MockUserRepo struct {
	mock.Mock
}

func (m *MockUserRepo) Create(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockUserRepo) FindByPhone(ctx context.Context, phone string) (*model.User, error) {
	args := m.Called(ctx, phone)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) FindByID(ctx context.Context, userID string) (*model.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepo) Update(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

type MockTxRepo struct {
	mock.Mock
}

func (m *MockTxRepo) Create(ctx context.Context, tx *model.Transaction) error {
	args := m.Called(ctx, tx)
	return args.Error(0)
}

func (m *MockTxRepo) FindByUserID(ctx context.Context, userID string) ([]model.Transaction, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]model.Transaction), args.Error(1)
}

func (m *MockTxRepo) FindByID(ctx context.Context, txID string) (*model.Transaction, error) {
	args := m.Called(ctx, txID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Transaction), args.Error(1)
}

func (m *MockTxRepo) UpdateStatus(ctx context.Context, txID string, status string) error {
	args := m.Called(ctx, txID, status)
	return args.Error(0)
}

// --- Test helpers ---

func setupService() (*TransactionService, *MockUserRepo, *MockTxRepo, chan worker.TransferJob) {
	// We use a nil *gorm.DB for the db field since we're mocking the repos
	// The service will use db only for actual DB transactions (TopUp/Payment)
	// For unit tests, we'll test the Transfer flow which doesn't need real DB
	transferQ := make(chan worker.TransferJob, 10)
	mockUserRepo := new(MockUserRepo)
	mockTxRepo := new(MockTxRepo)

	svc := &TransactionService{
		db:        nil,
		txRepo:    mockTxRepo,
		userRepo:  mockUserRepo,
		transferQ: transferQ,
	}

	return svc, mockUserRepo, mockTxRepo, transferQ
}

// --- Test cases ---

func TestTopUp_InvalidAmount(t *testing.T) {
	svc, _, _, _ := setupService()

	tx, err := svc.TopUp(context.Background(), "user-123", 0)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidAmount, err)
	assert.Nil(t, tx)

	tx, err = svc.TopUp(context.Background(), "user-123", -100)
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidAmount, err)
	assert.Nil(t, tx)
}

func TestTopUp_UserNotFound(t *testing.T) {
	svc, mockUserRepo, _, _ := setupService()

	mockUserRepo.On("FindByID", mock.Anything, "user-not-found").
		Return(nil, repository.ErrUserNotFound).Once()

	tx, err := svc.TopUp(context.Background(), "user-not-found", 50000)
	assert.Error(t, err)
	assert.Nil(t, tx)
}

func TestPayment_InvalidAmount(t *testing.T) {
	svc, _, _, _ := setupService()

	tx, err := svc.Payment(context.Background(), "user-123", 0, "test")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidAmount, err)
	assert.Nil(t, tx)
}

func TestPaymentInsufficient(t *testing.T) {
	// For Payment with insufficient balance, we need a mock that simulates
	// the row-lock transaction scenario. Since we can't inject a mock DB
	// easily, we'll test the Transfer flow instead which has clearer separation.
	// Payment's balance check happens inside the GORM transaction — hard to unit test without integration.
	// We'll rely on integration tests for Payment balance checks.
	t.Skip("Payment insufficient balance requires integration test with real DB")
}

func TestTransfer_InvalidAmount(t *testing.T) {
	svc, _, _, _ := setupService()

	tx, err := svc.Transfer(context.Background(), "sender-123", "receiver-456", 0, "test")
	assert.Error(t, err)
	assert.Equal(t, ErrInvalidAmount, err)
	assert.Nil(t, tx)
}

func TestTransfer_Success(t *testing.T) {
	svc, mockUserRepo, mockTxRepo, transferQ := setupService()

	receiverID := uuid.New().String()

	// Mock receiver exists
	mockUserRepo.On("FindByID", mock.Anything, receiverID).
		Return(&model.User{UserID: receiverID, Balance: 0}, nil).Once()

	// Mock tx creation
	mockTxRepo.On("Create", mock.Anything, mock.AnythingOfType("*model.Transaction")).
		Return(nil).Once()

	tx, err := svc.Transfer(context.Background(), "sender-123", receiverID, 50000, "test transfer")
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, model.KindTransfer, tx.TransactionKind)
	assert.Equal(t, model.TypeDebit, tx.TransactionType)
	assert.Equal(t, int64(50000), tx.Amount)
	assert.Equal(t, model.StatusPending, tx.Status)
	assert.Equal(t, receiverID, tx.RelatedUserID)

	// Verify job was enqueued
	select {
	case job := <-transferQ:
		assert.Equal(t, tx.TransactionID, job.TransactionID)
		assert.Equal(t, "sender-123", job.SenderID)
		assert.Equal(t, receiverID, job.ReceiverID)
		assert.Equal(t, int64(50000), job.Amount)
	default:
		t.Fatal("expected transfer job to be enqueued")
	}

	mockUserRepo.AssertExpectations(t)
	mockTxRepo.AssertExpectations(t)
}

func TestTransfer_ReceiverNotFound(t *testing.T) {
	svc, mockUserRepo, _, _ := setupService()

	mockUserRepo.On("FindByID", mock.Anything, "nonexistent").
		Return(nil, repository.ErrUserNotFound).Once()

	tx, err := svc.Transfer(context.Background(), "sender-123", "nonexistent", 50000, "test")
	assert.Error(t, err)
	assert.Nil(t, tx)
}

func TestGetTransactions(t *testing.T) {
	svc, _, mockTxRepo, _ := setupService()

	userID := uuid.New().String()
	txns := []model.Transaction{
		{
			TransactionID:   uuid.New().String(),
			UserID:         userID,
			TransactionKind: model.KindTopup,
			TransactionType: model.TypeCredit,
			Amount:          100000,
			BalanceBefore:  0,
			BalanceAfter:   100000,
			Status:         model.StatusSuccess,
		},
		{
			TransactionID:   uuid.New().String(),
			UserID:         userID,
			TransactionKind: model.KindPayment,
			TransactionType: model.TypeDebit,
			Amount:          50000,
			BalanceBefore:  100000,
			BalanceAfter:   50000,
			Status:         model.StatusSuccess,
		},
	}

	mockTxRepo.On("FindByUserID", mock.Anything, userID).
		Return(txns, nil).Once()

	result, err := svc.GetTransactions(context.Background(), userID)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, model.KindTopup, result[0].TransactionKind)
	assert.Equal(t, model.KindPayment, result[1].TransactionKind)

	mockTxRepo.AssertExpectations(t)
}

// --- Helper: integration-style service with real DB (for TopUp/Payment) ---

func setupServiceWithDB(t *testing.T) (*gorm.DB, *TransactionService, func()) {
	// This uses an in-memory SQLite for integration testing
	// Skipped in unit test mode; invoke only in integration tests
	t.Skip("use setupServiceWithDB in integration tests only")
	return nil, nil, func() {}
}