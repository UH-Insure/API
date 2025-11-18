# SAW + Cryptol + C Compiler API

This is a containerized API that:

Compiles & runs Cryptol  
Compiles & runs SAW proofs  
Compiles C code to LLVM bytecode for SAW  
Stores all submitted code in SQLite  
Allows password-protected API key access  
Allows Google Colab clients to evaluate models

---

## Endpoints

### POST `/run/cryptol`
Run Cryptol code.

### POST `/run/saw`
Run SAW scripts.  
If script references C files, they are compiled automatically.

### POST `/upload`
Upload cryptol/saw/c files (stored in SQLite + disk).

### GET `/files/{id}`
Download stored file.

---

## Environment Variables
API key will be printed on startup
API_KEY=yourapikey

podman build -t crysawapi .
podman run -p 8443:8443 crysawapi