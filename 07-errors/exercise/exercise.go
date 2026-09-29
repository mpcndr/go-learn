//go:build !solution

package errs

// ParseAge แปลง string เป็นอายุ
//   - แปลงเป็นตัวเลขไม่ได้ → &ValidationError{Field: "age", Reason: "not a number"}
//   - น้อยกว่า 0 หรือมากกว่า 150 → &ValidationError{Field: "age", Reason: "out of range"}
//
// คำใบ้: strconv.Atoi
func ParseAge(s string) (int, error) {
	// TODO
	return 0, nil
}

// FindUser หา user จาก map users (ใน types.go) ถ้าไม่พบคืน ErrNotFound
func FindUser(id int) (User, error) {
	// TODO
	return User{}, nil
}

// LoadProfile เรียก FindUser แล้วคืนข้อความ "<Name> <<Email>>" เช่น "Ann <ann@example.com>"
// ถ้า FindUser พัง ให้ wrap error ด้วยบริบท: "load profile <id>: <error เดิม>"
// เช่น "load profile 99: not found" และ errors.Is(err, ErrNotFound) ต้องยังเป็น true
func LoadProfile(id int) (string, error) {
	// TODO: fmt.Errorf กับ %w
	return "", nil
}

// ValidateAll ตรวจ User ทุกข้อ แล้วคืน error ทั้งหมดรวมกัน (ไม่ใช่หยุดที่ข้อแรก)
//   - Name ว่าง → &ValidationError{Field: "name", Reason: "required"}
//   - Email ไม่มี "@" → &ValidationError{Field: "email", Reason: "must contain @"}
//
// ถ้าไม่มีข้อผิดพลาดเลยต้องคืน nil
// คำใบ้: errors.Join (สังเกต: errors.Join() ที่ไม่มี argument ที่ไม่ใช่ nil จะคืน nil ให้เอง)
func ValidateAll(u User) error {
	// TODO
	return nil
}

// SafeDivide หาร a/b โดยใช้การหาร int ธรรมดา (ซึ่งจะ panic ถ้า b == 0)
// ห้ามเช็ค b == 0 ตรง ๆ — ให้ใช้ defer + recover แปลง panic เป็น error แทน
// (ฝึก mechanism เท่านั้น ในงานจริงให้เช็คก่อนหารจะดีกว่า)
func SafeDivide(a, b int) (result int, err error) {
	// TODO
	return a / b, nil
}
