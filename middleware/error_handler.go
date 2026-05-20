package middleware

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler เป็น Middleware สำหรับดักจับและรวมการตอบกลับ Error ไว้ที่เดียว
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		// ให้ Request วิ่งไปทำงานที่ Handler ปกติก่อน
		c.Next()

		// หลังจาก Handler ทำงานเสร็จ จะกลับมาตรวจสอบว่ามี Error เกิดขึ้นใน Context ไหม
		if len(c.Errors) > 0 {
			// ดึง Error ออกมา (สามารถวนลูปดูได้หากมีหลาย Error)
			err := c.Errors.Last()
			log.Printf("App Error: %v\n", err.Err)

			// ตรวจสอบว่าใน Handler ได้กำหนด Status Code ไว้ล่วงหน้าไหม (ถ้ายังจะเป็น 200)
			status := c.Writer.Status()
			if status == http.StatusOK {
				status = http.StatusInternalServerError // ค่าเริ่มต้นถ้าเกิด Error แต่ไม่ได้ตั้ง Status
			}

			// ส่ง JSON Response ให้ Client แบบมาตรฐาน
			c.JSON(status, gin.H{
				"error":   true,
				"message": err.Error(),
			})
		}
	}
}
