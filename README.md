# Penugasan-Netics-CI-CD


# Laporan Penugasan Modul CI/CD  
**Open Recruitment NETICS 2025**  
**Author**: Iftala Zahri Sukmana  

---

## 1. Deskripsi Singkat  
Tugas ini mengimplementasikan modul CI/CD untuk sebuah API sederhana dengan endpoint `/health`. API menampilkan:
```json
{
  "nama": "Rogelio Kenny Arisandi",
  "nrp": "5025231074",
  "status": "UP",
  "timestamp": 1680301234,
  "uptime": 120
}```

Seluruh proses—mulai dari build, containerization, hingga deployment ke VPS—diotomasi menggunakan GitHub Actions.

# 2. Teknologi yang Digunakan
	•	Go (net/http): Bahasa dan library standar untuk membuat HTTP server ringan.
	•	Docker: Memisahkan tahap build (dengan image Go) dan final (Alpine minimal) untuk image yang kecil.
	•	GitHub Actions: Platform CI/CD untuk otomatisasi build & deploy.
	•	Docker Hub: Registry publik untuk menyimpan image (enyytime/health-api:latest).
	•	AWS EC2 (Ubuntu 20.04, t2.micro): VPS tempat menjalankan container.

# 3. Implementasi API (main.go)
```go
// Handler menghitung uptime sejak server startTime dan mengembalikan JSON.
func healthHandler(w http.ResponseWriter, r *http.Request) {
  now := time.Now()
  uptime := int64(now.Sub(startTime).Seconds())
  resp := HealthResponse{
    Nama:      "...", NRP: "...", Status: "UP",
    Timestamp: now.Unix(), Uptime: uptime,
  }
  w.Header().Set("Content-Type", "application/json")
  json.NewEncoder(w).Encode(resp)
}
```
startTime dicatat sekali saat main() dijalankan.
Timestamp dan Uptime dihitung dinamis setiap request.

# 4. Membuar Docker Image dan Dockerfile, lalu push ke docker hub
```Dockerfile
# Stage 1: build binary statis
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o health-api .

# Stage 2: final image 
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/health-api .
EXPOSE 8080
ENTRYPOINT ["./health-api"]
```
Setelah

## Push image ke docker hub

Sebelum deployment, kita harus mengunggah image ke Docker Hub:
1. Login ke Docker Hub (sekali di mesin lokal):
2. Tag image lokal
```
docker tag go-health-api enyytime/health-api:latest
```
3. Push ke registry
```
docker push enyytime/health-api:latest
```
 

# 6. Deployment di AWS EC2

	1.	Provision EC2 (Ubuntu 20.04, t2.micro).
	2.	Security Group: buka port 22 (SSH) dan 80 (HTTP).
	3.	Install Docker dan tambahkan user ke grup docker agar bisa tanpa sudo:

### Berikut adalah Command yang di jalankan untuk pull image docker

1. install docker di VPS
```
sudo apt update
sudo apt install docker.io -y
sudo usermod -aG docker $USER
```
2. pull image dari docker hub kita
```
docker pull enyytime/health-api:latest
docker run -d --name health-api -p 80:8080 enyytime/health-api:latest
```


