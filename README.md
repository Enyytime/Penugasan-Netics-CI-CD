# Penugasan-Netics-CI-CD


# Laporan Penugasan Modul CI/CD  
**Open Recruitment NETICS 2025**  
**Name**: Rogelio Kenny Arisandi
**NRP**: 5025231074

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
}
```

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

1. Melakukan koneksi SSH dengan VPS
- Copy IP ini untuk melakukan SSH
![image](https://github.com/user-attachments/assets/4fbc4377-4394-4ea6-9925-74966ee44b0c)
- Lakukan command ini
![image](https://github.com/user-attachments/assets/a85ce7e1-9741-4647-88c0-770f9e90d2cf)


2. install docker di VPS
```
sudo apt update
sudo apt install docker.io -y
sudo usermod -aG docker $USER
```
3. pull image dari docker hub kita
```
docker pull enyytime/health-api:latest
docker run -d --name health-api -p 80:8080 enyytime/health-api:latest
```

# 7. CI/CD
- Buat directory `.github/workflows` lalu buat file `deploy.yml`
```yml
name: CI/CD Pipeline

on:
  push:
    branches:
      - main

jobs:
  build:
    name: Build and Push Docker Image
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Login to DockerHub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Build and Push Docker Image
        uses: docker/build-push-action@v4
        with:
          context: .
          file: Dockerfile
          push: true
          tags: enyytime/health-api:latest

  deploy:
    name: Deploy to VPS
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: SSH Deploy to VPS
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.VPS_HOST }}
          username: ${{ secrets.VPS_USER }}
          key: ${{ secrets.VPS_SSH_KEY }}
          script: |
            sudo docker pull enyytime/health-api:latest
            sudo docker stop health-api || true
            sudo docker rm health-api || true
            sudo docker run -d --name health-api -p 80:8080 enyytime/health-api:latest
```


name: CI/CD Pipeline

on:
  push:
    branches:
      - main

jobs:
  build:
    name: Build and Push Docker Image
    runs-on: ubuntu-latest
    steps:
      - name: Checkout Code
        uses: actions/checkout@v3

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      - name: Login to DockerHub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKER_USERNAME }}
          password: ${{ secrets.DOCKER_PASSWORD }}

      - name: Build and Push Docker Image
        uses: docker/build-push-action@v4
        with:
          context: .
          file: Dockerfile
          push: true
          tags: enyytime/health-api:latest

  deploy:
    name: Deploy to VPS
    needs: build
    runs-on: ubuntu-latest
    steps:
      - name: SSH Deploy to VPS
        uses: appleboy/ssh-action@master
        with:
          host: ${{ secrets.VPS_HOST }}
          username: ${{ secrets.VPS_USER }}
          key: ${{ secrets.VPS_SSH_KEY }}
          script: |
            sudo docker pull enyytime/health-api:latest
            sudo docker stop health-api || true
            sudo docker rm health-api || true
            sudo docker run -d --name health-api -p 80:8080 enyytime/health-api:latest

### Penjelasan Kode

#### 1. Build dan Push Docker Image ke Docker Hub

Pada tahap ini, GitHub Actions akan:

- Mengambil kode terbaru dari repository (actions/checkout@v3).
- Menyiapkan Docker Buildx untuk membangun image.
- Melakukan login ke Docker Hub menggunakan secrets yang telah disimpan di repository.
- - Membangun image berdasarkan Dockerfile dan langsung mem-push ke Docker Hub.

#### 2. Deploy ke VPS menggunakan SSH
Pada tahap ini, setelah Docker image berhasil dipush ke Docker Hub, workflow akan:
Menjalankan perintah:
- Menarik (pull) image terbaru dari Docker Hub.

- Menghentikan container lama (docker stop), jika ada.

- Menghapus container lama (docker rm), jika ada.

- Menjalankan container baru dengan port mapping 80:8080.

### Isi Secrets di githuub actions

DOCKER_USERNAME: Docker Hub username.

DOCKER_PASSWORD: Docker Hub password.

VPS_HOST: VPS public IP.

VPS_USER: ubuntu.

VPS_SSH_KEY: SSH private key (.pem in my case)

# 8. Tes CI/CD
Cara tes nya tinggal melakukan update lalu kita cek di github action


![image](https://github.com/user-attachments/assets/a43a25e9-4bce-4615-a74e-beea25cd0051)



![image](https://github.com/user-attachments/assets/e542aaa6-5e23-48a0-ba5d-aff9a3acfaf4)

Penjelasan:
Uptime nya pada saat `1743948946`, artinya ini di tanggal 6 april 14:15 UTC, ditambah 7 karena kita +7 jadinya jam 9. (nanya chatgpt karena gak bisa baca itu angkanya gimana)
