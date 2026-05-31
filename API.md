# go-imgbed API Reference

## Base URL
API 请求地址与上传页面相同，例如：

```
https://example.com/upload
```

所有操作均通过 **POST** 请求发送，并使用自定义请求头 `X-Type` 指定操作类型。

---

# Authentication

所有请求都需要 **HTTP Basic Authentication**。

如果认证失败，服务器返回：

```
HTTP/1.1 401 Unauthorized
WWW-Authenticate: Basic realm="restricted"
```

---

# Response Format

所有响应均为 **JSON 格式**。

### Success Example

```json
{
  "state": "/file/xxxxx.jpg",
  "ok": true
}
```

### Failure Example

```json
{
  "state": "error message",
  "ok": false
}
```

---

# 1. Upload Image

上传图片到服务器。

## Request

**Method**

```
POST
```

**Headers**

```
X-Type: upload
```

**Body**

```
multipart/form-data
```

**Fields**

```
file : image file
```

## Constraints

- 支持格式：`jpg` `png` `webp`
- 最大大小：`15 MB`
- 服务端文件名由 **SHA256 Hash** 自动生成

## Example

```bash
curl -u user:pass \
  -H "X-Type: upload" \
  -F "file=@test.jpg" \
  https://example.com/upload
```

## Success Response

```json
{
  "state": "/file/Q2xK...abc.jpg",
  "ok": true
}
```

`state` 字段包含文件的公开访问路径。

### Full URL Example

```
https://example.com/file/Q2xK...abc.jpg
```

## Failure Examples

```json
{
  "state": "only jpg, png, webp allowed",
  "ok": false
}
```

```json
{
  "state": "too big",
  "ok": false
}
```

---

# 2. List Images

获取已上传图片列表。

## Request

**Method**

```
POST
```

**Headers**

```
X-Type: list
```

## Example

```bash
curl -u user:pass \
  -H "X-Type: list" \
  -X POST \
  https://example.com/upload
```

## Response

```json
{
  "abc123.jpg": 0,
  "def456.png": 0,
  "xyz789.webp": 0
}
```

## Notes

- `key` 为文件名
- `value` 恒为 `0`（占位符）

### Image URL Format

```
https://example.com/file/{filename}
```

---

# 3. Delete Image

删除已上传图片。

## Request

**Method**

```
POST
```

**Headers**

```
X-Type: del
Content-Type: application/json
```

**Body**

```json
{
  "file": "abc123.jpg"
}
```

## Example

```bash
curl -u user:pass \
  -H "X-Type: del" \
  -H "Content-Type: application/json" \
  -d '{"file":"abc123.jpg"}' \
  https://example.com/upload
```

## Success Response

```json
{
  "state": "delete OK",
  "ok": true
}
```

## Failure Response

```json
{
  "state": "file not found",
  "ok": false
}
```

---

# 4. Error Responses

## Unauthorized

```
HTTP/1.1 401 Unauthorized
```

## Unknown Action

```json
{
  "state": "type not found",
  "ok": false
}
```

---

# 5. Public File Access

上传的文件可以通过以下路径访问：

```
/file/{filename}
```

### Example

```
https://example.com/file/Q2xK...abc.jpg
```