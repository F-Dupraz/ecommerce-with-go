package dto

import (
  "regexp"

  "github.com/go-playground/validator/v10"
)

func ValidatePassword(fl validator.FieldLevel) bool {
  password := fl.Field().String()

  hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(password)
  hasLower := regexp.MustCompile(`[a-z]`).MatchString(password)
  hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
  hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>\/?]`).MatchString(password)

  return hasUpper && hasLower && hasNumber && hasSpecial
}

func ValidateISO3166Alpha2(fl validator.FieldLevel) bool {
  country := fl.Field().String()

  validCountries := map[string]bool{
	"AD": true, "AE": true, "AF": true, "AG": true, "AI": true, "AL": true, "AM": true,
	"AO": true, "AQ": true, "AR": true, "AS": true, "AT": true, "AU": true, "AW": true,
	"AX": true, "AZ": true, "BA": true, "BB": true, "BD": true, "BE": true, "BF": true,
	"BG": true, "BH": true, "BI": true, "BJ": true, "BL": true, "BM": true, "BN": true,
	"BO": true, "BQ": true, "BR": true, "BS": true, "BT": true, "BV": true, "BW": true,
	"BY": true, "BZ": true, "CA": true, "CC": true, "CD": true, "CF": true, "CG": true,
	"CH": true, "CI": true, "CK": true, "CL": true, "CM": true, "CN": true, "CO": true,
	"CR": true, "CU": true, "CV": true, "CW": true, "CX": true, "CY": true, "CZ": true,
	"DE": true, "DJ": true, "DK": true, "DM": true, "DO": true, "DZ": true, "EC": true,
	"EE": true, "EG": true, "EH": true, "ER": true, "ES": true, "ET": true, "FI": true,
	"FJ": true, "FK": true, "FM": true, "FO": true, "FR": true, "GA": true, "GB": true,
	"GD": true, "GE": true, "GF": true, "GG": true, "GH": true, "GI": true, "GL": true,
	"GM": true, "GN": true, "GP": true, "GQ": true, "GR": true, "GS": true, "GT": true,
	"GU": true, "GW": true, "GY": true, "HK": true, "HM": true, "HN": true, "HR": true,
	"HT": true, "HU": true, "ID": true, "IE": true, "IL": true, "IM": true, "IN": true,
	"IO": true, "IQ": true, "IR": true, "IS": true, "IT": true, "JE": true, "JM": true,
	"JO": true, "JP": true, "KE": true, "KG": true, "KH": true, "KI": true, "KM": true,
	"KN": true, "KP": true, "KR": true, "KW": true, "KY": true, "KZ": true, "LA": true,
	"LB": true, "LC": true, "LI": true, "LK": true, "LR": true, "LS": true, "LT": true,
	"LU": true, "LV": true, "LY": true, "MA": true, "MC": true, "MD": true, "ME": true,
	"MF": true, "MG": true, "MH": true, "MK": true, "ML": true, "MM": true, "MN": true,
	"MO": true, "MP": true, "MQ": true, "MR": true, "MS": true, "MT": true, "MU": true,
	"MV": true, "MW": true, "MX": true, "MY": true, "MZ": true, "NA": true, "NC": true,
	"NE": true, "NF": true, "NG": true, "NI": true, "NL": true, "NO": true, "NP": true,
	"NR": true, "NU": true, "NZ": true, "OM": true, "PA": true, "PE": true, "PF": true,
	"PG": true, "PH": true, "PK": true, "PL": true, "PM": true, "PN": true, "PR": true,
	"PS": true, "PT": true, "PW": true, "PY": true, "QA": true, "RE": true, "RO": true,
	"RS": true, "RU": true, "RW": true, "SA": true, "SB": true, "SC": true, "SD": true,
	"SE": true, "SG": true, "SH": true, "SI": true, "SJ": true, "SK": true, "SL": true,
	"SM": true, "SN": true, "SO": true, "SR": true, "SS": true, "ST": true, "SV": true,
	"SX": true, "SY": true, "SZ": true, "TC": true, "TD": true, "TF": true, "TG": true,
	"TH": true, "TJ": true, "TK": true, "TL": true, "TM": true, "TN": true, "TO": true,
	"TR": true, "TT": true, "TV": true, "TW": true, "TZ": true, "UA": true, "UG": true,
	"UM": true, "US": true, "UY": true, "UZ": true, "VA": true, "VC": true, "VE": true,
	"VG": true, "VI": true, "VN": true, "VU": true, "WF": true, "WS": true, "YE": true,
	"YT": true, "ZA": true, "ZM": true, "ZW": true,
  }

  return validCountries[country]
}

func RegisterCustomValidators(v *validator.Validate) error {
  if err := v.RegisterValidation("password", ValidatePassword); err != nil {
	return err
  }
	
  if err := v.RegisterValidation("iso3166_1_alpha2", ValidateISO3166Alpha2); err != nil {
	return err
  }

  return nil
}

func NewValidator() (*validator.Validate, error) {
  v := validator.New()

  if err := RegisterCustomValidators(v); err != nil {
	return nil, err
  }

  return v, nil
}

func FormatValidationErrors(err error) map[string]string {
  errors := make(map[string]string)

  if validationErrors, ok := err.(validator.ValidationErrors); ok {
	for _, e := range validationErrors {
	  field := e.Field()
	  tag := e.Tag()

	  switch tag {
	  case "required":
		errors[field] = field + " is required"
	  case "email":
		errors[field] = field + " must be a valid email"
	  case "min":
		errors[field] = field + " must be at least " + e.Param() + " characters"
	  case "max":
		errors[field] = field + " must be at most " + e.Param() + " characters"
	  case "alphanum":
		errors[field] = field + " must contain only letters and numbers"
	  case "password":
		errors[field] = "Password must contain at least 1 uppercase, 1 lowercase, 1 number, and 1 special character"
	  case "iso3166_1_alpha2":
		errors[field] = field + " must be a valid ISO 3166-1 alpha-2 country code"
	  case "uuid":
		errors[field] = field + " must be a valid UUID"
	  default:
		errors[field] = field + " failed " + tag + " validation"
	  }
	}
  }

  return errors
}

