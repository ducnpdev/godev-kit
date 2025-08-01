# VietQR Integration

The VietQR module provides APIs to generate, inquire, and update the status of VietQR codes.

## 1. Generate QR Code

- **Endpoint:** `POST /v1/vietqr/gen`
- **Description:** Generate a new VietQR code.
- **Request Body:**
  ```json
  {
    "accountNo": "string",      // Required. Bank account number.
    "amount": "string",         // Required. Amount to be paid.
    "description": "string",    // Optional. Payment description.
    "mcc": "string",            // Optional. Merchant Category Code.
    "receiverName": "string"    // Optional. Name of the receiver.
  }
  ```
- **Response:**
  ```json
  {
    "id": "string",             // QR code ID.
    "status": "string",         // Status: generated, in-process, paid, fail, timeout.
    "content": "string"         // QR code content (e.g., base64 or URL).
  }
  ```
- **Success Code:** 200

---

## 2. Inquiry QR Status

- **Endpoint:** `GET /v1/vietqr/inquiry/{id}`
- **Description:** Get the status and content of a VietQR code by its ID.
- **Path Parameter:**
  - `id` (string): QR code ID.
- **Response:**
  ```json
  {
    "id": "string",
    "status": "string",         // Status: generated, in-process, paid, fail, timeout.
    "content": "string"
  }
  ```
- **Success Code:** 200

---

## 3. Update QR Status

- **Endpoint:** `PUT /v1/vietqr/update/{id}`
- **Description:** Update the status of a VietQR code.
- **Path Parameter:**
  - `id` (string): QR code ID.
- **Request Body:**
  ```json
  {
    "status": "string"          // Required. One of: in-process, paid, fail, timeout.
  }
  ```
- **Response:**
  ```json
  {
    "status": "ok"
  }
  ```
- **Success Code:** 200

---

### Status values
- `generated`
- `in-process`
- `paid`
- `fail`
- `timeout` 


## Diagram
```
flowchart TD
    A["Client<br/>POST /v1/vietqr/gen"] --> B["HTTP Controller<br/>generateQR"]
    B --> C["UseCase<br/>GenerateQR"]
    C --> D["External API Repo<br/>GenerateQR (vietqr lib)"]
    D --> E["VietQR Content"]
    C --> F["Persistent Repo<br/>Store VietQR in DB"]
    F --> G["DB: vietqr table"]
    C --> H["Return VietQR Entity (id, status, content)"]
    H --> I["Client Receives QR Info"]
    
    subgraph Status Inquiry/Update
      J["Client<br/>GET /v1/vietqr/inquiry/{id}"] --> K["HTTP Controller<br/>inquiryQR"]
      K --> L["UseCase<br/>InquiryQR"]
      L --> M["Persistent Repo<br/>FindByID"]
      M --> N["DB: vietqr table"]
      L --> O["Return VietQR Entity"]
      O --> P["Client Receives Status"]
      
      Q["Client<br/>PUT /v1/vietqr/update/{id}"] --> R["HTTP Controller<br/>updateStatus"]
      R --> S["UseCase<br/>UpdateStatus"]
      S --> T["Persistent Repo<br/>UpdateStatus"]
      T --> U["DB: vietqr table"]
      S --> V["Return Success"]
      V --> W["Client Receives Update Ack"]
    end
    style Status\ Inquiry\/Update fill:#f9f,stroke:#333,stroke-width:2
```