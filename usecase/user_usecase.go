package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"gym-app/domain"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo domain.UserRepository
	validate *validator.Validate
}

func NewUserUsecase(ur domain.UserRepository, val *validator.Validate) domain.UserUsecase {
	return &userUsecase{
		userRepo: ur,
		validate: val,
	}
}

func (u *userUsecase) Register(req domain.UserRegisterRequest) (*domain.User, error) {
	// Validasi input
	if err := u.validate.Struct(req); err != nil {
		return nil, errors.New("semua field wajib diisi dan format email harus valid")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("gagal memproses enkripsi password")
	}

	// Mapping ke struct User
	user := &domain.User{
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		Password: string(hashedPassword),
		Weight:   req.Weight,
		Height:   req.Height,
	}

	// Simpan ke database
	err = u.userRepo.Create(user)
	if err != nil {
		return nil, errors.New("gagal mendaftar, email atau username mungkin sudah digunakan")
	}

	return user, nil
}

// Skeleton untuk fungsi selanjutnya
func (u *userUsecase) Login(email, password string) (string, error) {
	// 1. Cari user di database berdasarkan email
	user, err := u.userRepo.GetByEmail(email)
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	// 2. Bandingkan password dari request dengan password hash di database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	// 3. Jika cocok, buat klaim (payload) untuk JWT Token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token berlaku 72 jam
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 4. Tanda tangani token menggunakan JWT_SECRET dari file .env
	secret := viper.GetString("JWT_SECRET")
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("terjadi kesalahan saat membuat token")
	}

	return signedToken, nil
}

func (u *userUsecase) GetProfile(id int) (*domain.User, float64, error) {
	// 1. Ambil data user dari database
	user, err := u.userRepo.GetByID(id)
	if err != nil {
		return nil, 0, errors.New("user tidak ditemukan")
	}

	heightInMeter := float64(user.Height) / 100.0
	url := fmt.Sprintf("https://body-mass-index-bmi-calculator.p.rapidapi.com/metric?weight=%d&height=%.2f", user.Weight, heightInMeter)

	fmt.Println("--- DEBUGGING API ---")
	fmt.Println("1. URL yang ditembak:", url)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("X-RapidAPI-Key", viper.GetString("RAPIDAPI_KEY"))
	req.Header.Add("X-RapidAPI-Host", "body-mass-index-bmi-calculator.p.rapidapi.com")

	// Eksekusi HTTP Request
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("2. ERROR: Gagal menghubungi RapidAPI:", err.Error())
		// Kembalikan error sungguhan agar handler tahu prosesnya gagal
		return user, 0, errors.New("gagal menghubungi layanan BMI pihak ketiga")
	}
	defer res.Body.Close()

	// Baca response dari RapidAPI
	body, _ := io.ReadAll(res.Body)

	fmt.Println("2. Status Code RapidAPI:", res.StatusCode)
	fmt.Println("3. Response Asli RapidAPI:", string(body))
	fmt.Println("---------------------")

	if res.StatusCode != 200 {
		// Beritahu handler bahwa API menolak request kita
		return user, 0, fmt.Errorf("layanan BMI mengembalikan status error: %d", res.StatusCode)
	}

	var bmiResponse struct {
		Bmi float64 `json:"bmi"`
	}

	err = json.Unmarshal(body, &bmiResponse)
	if err != nil {
		fmt.Println("4. ERROR Unmarshal JSON:", err.Error())
		// Hentikan proses dan kembalikan error
		return user, 0, errors.New("gagal memproses data balasan dari layanan BMI")
	}

	return user, bmiResponse.Bmi, nil
}
