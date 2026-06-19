# 🎯 โจทย์: Centralized CI/CD Pipeline on GitHub Actions (Go & PostgreSQL)

**สถานะปัจจุบัน:** เรามี Source Code ของระบบ Task Tracker API (Go + PostgreSQL) พร้อมไฟล์ Migration และ Unit Test เรียบร้อยแล้ว

**ภารกิจของคุณ:** คือการสร้างท่อส่งของอัตโนมัติ (Pipeline) โดยย้ายกระบวนการควบคุมทั้ง CI และ CD มาไว้ที่ **GitHub Actions** ทั้งหมด เพื่อสั่งงานไปยัง **Render** 

> **เงื่อนไขสำคัญ:** "PostgreSQL จะถูก Deploy แค่ครั้งเดียวตอนเริ่มโปรเจกต์ ส่วนแอป Go จะถูก Deploy ใหม่ทุกครั้งที่มีการอัปเดตโค้ด"

---

## 🏗️ การแบ่งสถาปัตยกรรม Pipeline (Requirements)

ให้ดีไซน์ Workflow บน GitHub Actions โดยแยกการทำงานออกเป็น 2 ส่วนชัดเจน ดังนี้:

### 1. Pipeline ฝั่ง Database (PostgreSQL) ➔ ทำครั้งเดียว

* **การทำงาน:** ให้เขียน Workflow สำหรับ Setup ตัว PostgreSQL บน Render ให้พร้อมใช้งาน
* **เงื่อนไข:** Pipeline นี้จะถูกรันแบบ Manual (ใช้ `workflow_dispatch`) หรือรันแค่ครั้งแรกครั้งเดียวเท่านั้น เพื่อจองพื้นที่และเปิด Service Database บน Cloud โดยจะไม่ยุ่งกับมันอีกในการส่งโค้ดรอบถัดๆ ไป

### 2. Pipeline ฝั่งแอปพลิเคชัน (Go Backend) ➔ ทำทุกครั้งที่ Push `dev`

ให้เขียน Workflow กำหนดเงื่อนไขให้ทำงานอัตโนมัติ **"ทุกครั้งที่มีการ `git push` เข้าหา Branch `dev`"** โดยต้องไล่สเต็ปดังนี้:

* **[สเต็ป CI]:** รัน Code Linter และ Automated Testing (Unit Test) เพื่อตรวจความถูกต้องของโค้ด Go
* **[สเต็ป Containerize]:** บิวด์โค้ดให้เป็น Docker Image ด้วยเทคนิค Multi-stage Build ผ่าน `Dockerfile`
* **[สเต็ป CD]:** ส่งคำสั่งและอิมเมจที่บิวด์เสร็จแล้วไป Deploy ที่ Web Service บน Render โดยตรงผ่าน Render Deploy Hook หรือ GitHub Actions Integration *(ห้ามเปิด Auto-Deploy ฝั่งหน้าเว็บ Render)*
* **[สเต็ป Auto-Migration]:** เมื่อแอป Go สตาร์ทบน Cloud มันจะต้องสั่งรันสคริปต์ Migration เพื่อไปสร้างหรืออัปเดตตารางข้อมูลใน PostgreSQL (ที่ถูกเปิดไว้จากข้อ 1) โดยอัตโนมัติ

---

## 🛑 ข้อห้ามและเงื่อนไขความปลอดภัย

1. **ห้ามเปิด GitHub Autodeploy บน Render:** การ Deploy ทั้งหมดต้องถูกทริกเกอร์และสั่งการมาจาก GitHub Actions เท่านั้น
2. **ห้าม Hardcode Credentials:** ค่าคอนฟิก เช่น โทเคนหรือ Deploy Hook ของ Render (`RENDER_TOKEN`, `RENDER_DEPLOY_HOOK_URL`) หรือรหัสผ่านฐานข้อมูล ห้ามใส่ไว้ในโค้ดหรือไฟล์ `.yml` เด็ดขาด ให้เก็บไว้ใน **GitHub Secrets** ของ Repository แล้วเรียกใช้ผ่านระบบ Environment Variable เท่านั้น

---

## 📦 สรุปผลลัพธ์สุดท้ายบนแดชบอร์ด Render (Definition of Done)

เมื่อระบบทำงานสำเร็จ หน้าจอ Render จะต้องแสดงผลดังนี้:

1. **Web Service (Go API):** ต้องแสดงสถานะว่าถูก Deploy ล่าสุดผ่านทาง GitHub Actions (ไม่ใช่จาก Render Build เอง) มี URL ที่ใช้งานได้จริง และมีตัวแปร `PORT` กับ `DATABASE_URL` ผูกอยู่
2. **Database Service (PostgreSQL):** เปิดใช้งานอยู่ปกติ มีเส้นประเชื่อมต่อ (Connection Line) วิ่งไปหาแอป Go และเมื่อเข้าไปดูในแท็บ Data จะต้องพบตารางข้อมูลที่ถูกเนรมิตขึ้นมาผ่านระบบ Auto-Migration เรียบร้อย
3. **GitHub Actions History:** ในหน้า GitHub ต้องแสดงสถานะผ่าน (สีเขียว ✅) ทั้งหมด โดยเห็นขั้นตอนการรัน Linter, Test, Build Docker และ Deploy ไป Render แสดงผลอย่างชัดเจน