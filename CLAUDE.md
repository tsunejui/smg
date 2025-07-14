# CLAUDE.md

## 📌 Project Name: Social Media Growth Engine

This project is a full-stack application integrating a web-based admin panel and a mobile app. The goal is to help users collect and repost content related to their favorite topics from various social media platforms, thereby growing their social following.

---

## 📦 Architecture and Development Guidelines

### Web Admin Panel:

* Framework: Next.js
* UI: Tailwind CSS
* Authentication: Google OAuth and Email/Password
* Database: PostgreSQL
* Migration Tool: go-migrate
* Communication: RPC between backend services written in Go
* Dev Environment: Docker Compose
* Repository Structure (Monorepo):

```
root/
├── apps/
│   ├── web/               # Next.js Admin Panel
│   └── app/               # Flutter Mobile App
├── cmd/                   # Various Go-based CLI tools for notifications, automation, etc.
├── pkg/                   # Shared Go modules
├── scripts/
├── Makefile
├── docker-compose.yml
└── docs/                  # Developer guides and documentation
```

### Mobile App:

* Framework: Flutter (iOS + Android support)
* Function: Content approval and reposting

### Backend Logic:

* Language: Go
* CLI Tools: Separated by features (e.g., push notifications, auto-repost scheduler)
* Communication: RPC among services

### Makefile Example:

```make
start:
	docker-compose up --build -d

stop:
	docker-compose down

web:
	cd apps/web && npm run dev

migrate:
	go run cmd/migrate/main.go
```

---

## 🗂 Admin Panel Functionalities

### 🔐 Authentication and Verification

* Google OAuth2 and Email/Password login
* Test account for local development
* Email verification flow
* Documentation: `docs/google-auth.md` provides instructions to set up GCP, enable APIs, create OAuth Client IDs, and configure callback URIs.

### 👤 User Dashboard

**Sidebar Navigation (中文):**

* 主題管理
* 媒體連結
* 文章管理
* 系統設定

**Navbar Menu (中文):**

* 個人頭像（修改個人資料）
* QR Code（供手機掃描登入）
* 登出

#### 主題管理

* 欄位：名稱、描述、關鍵字
* 可選擇多個目標平台（由管理者定義）
* 可綁定使用者自己的媒體帳號進行搜尋

#### 媒體連結

* 目前支援 X（Twitter）
* 使用者可連結自己的社群帳號
* 未來支援 Threads、Instagram、Facebook、TikTok 等

#### 文章管理

* 根據主題與媒體連結，擷取相關文章並分類
* 可轉貼文章至媒體帳號，並附上自訂描述或 AI 生成描述

#### 系統設定

* 設定自動轉貼邏輯：指定主題、頻率與平台
* 範例：每小時自動轉貼「科技新聞」到 X

### 🛠 管理者後台

**預設登入帳密：** `admin / admin`（首次登入後需強制修改密碼）

**Sidebar Navigation (中文):**

* 總覽
* 用戶管理
* 社群媒體設定
* 系統設定

**Navbar Menu (中文):**

* 個人頭像（修改密碼）
* 登出

#### 總覽

* 儀表板顯示活躍用戶與成長趨勢

#### 用戶管理

* 顯示所有用戶列表
* 點選彈出視窗可檢視與編輯用戶資料
* 查看其綁定的媒體帳號、主題、轉貼次數統計

#### 社群媒體設定

* 設定各平台的搜尋 API 或連線資訊
* 提供使用者在前台進行綁定或搜尋

#### 系統設定

* 設定 SMTP 寄信服務（例如新用戶註冊通知）

---

## 📱 Mobile App Functionalities

* 掃描 QR code 登入
* 查看依主題分類的文章
* 決定是否要轉貼內容至社群帳號

---

## 📘 Documentation (`/docs` folder)

* `google-auth.md`: Google OAuth2 setup guide
* `postgre-guideline.md`: PostgreSQL structure and naming conventions
* `golang-guideline.md`: Follow Uber's Go Style Guide
* `nextjs-guideline.md`: Next.js project layout and conventions
* `flutter-guideline.md`: Flutter code and UI best practices
* `rpc-spec.md`: RPC service definitions and format
* `dev-env.md`: Docker Compose dev setup instructions
* `makefile-commands.md`: Description of Makefile commands

---

## 🧠 Code Commentary and Guidelines

* All business logic must include inline comments
* Follow language-specific conventions:

  * Go: proper error handling, modular structure
  * Flutter: use of widget trees and stateless/stateful separation
  * Next.js: SSR/CSR usage where appropriate
* Formatting tools enforced:

  * `gofmt`, `prettier`, `flutter format`

---

## ❓Clarification and Additional Features

### Areas Needing Clarification

* **Repost Execution:** Should reposts be done via official APIs or through simulated automation?
* **AI Captioning:** Is GPT-style caption generation desired? If so, should it be run locally or via OpenAI API?
* **Post Approval Queue:** Who moderates content? Is it user-driven or admin-reviewed?
* **Media Crawling Logic:** Are we scraping content or consuming public APIs?

### Suggested Additional Features

* **Analytics Dashboard for Users:** Show how many reposts were made, follower changes over time
* **Webhook/IFTTT Integrations:** For advanced repost triggers or cross-platform syncing
* **Notification System:** Use Go CLI tool for scheduling reminders or alerting repost failures
* **Team Access for Admin Panel:** Admin roles (read-only, editor, full access)
* **Rate Limiting:** Prevent abuse by users over-posting to social media