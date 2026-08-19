---
project: go-split
topic: cach-chia-tien
language: vi
---

# Cách chia tiền

Đây là phần quan trọng nhất của go-split. Đọc kỹ để nhập khoản chi đúng, số dư mới đúng.

## Hai cách chia

### 1. Chia đều

Không gán số tiền riêng cho ai. Tổng hóa đơn chia đều cho tất cả người cùng chịu.

Ví dụ: 300.000 đồng, 3 người cùng chịu → mỗi người 100.000.

### 2. Chia tùy chỉnh kết hợp chia đều

Gán số tiền cố định cho một số người. Phần tiền còn lại chia đều cho những người **chưa** được gán số cố định.

Ví dụ: hóa đơn 300.000, ba người A, B, C.

- A chịu cố định 50.000 (ít hơn vì ăn ít)
- Còn 250.000 chia cho B và C → mỗi người 125.000

A không bị chia phần còn lại vì đã có số cố định.

## Người trả nhiều người

Nếu hai người cùng trả một hóa đơn, phần **đã trả** được tính chia đều cho những người trả.

Ví dụ hóa đơn 200.000, An và Bình cùng trả → mỗi người được tính đã trả 100.000.

Phần **phải chịu** vẫn theo danh sách người cùng chịu và cách chia ở trên, không phụ thuộc ai đưa tiền cho quán.

## Công thức dễ nhớ

Với mỗi người:

- **Đã trả** = phần người đó ứng (tổng hóa đơn chia đều cho những người trả, nếu họ có trong danh sách trả)
- **Phải chịu** = phần chia cho họ (cố định hoặc chia đều phần còn lại)
- **Số dư** = đã trả − phải chịu

Số dư dương: người khác đang nợ bạn.  
Số dư âm: bạn đang nợ người khác.  
Số dư 0: đã cân với khoản đó / nhóm đó.

## Việc cần nhớ khi nhập

- Phải có ít nhất một người trả và một người cùng chịu
- Người trả phải thuộc người cùng chịu
- Tổng các phần cố định không nên lớn hơn tổng hóa đơn; phần còn lại mới đủ chia cho người chưa gán
- Nếu gán cố định cho **mọi** người cùng chịu thì không còn ai để chia phần dư. Tránh trường hợp này, hoặc gán sao cho khớp tổng tiền
- Khi sửa người cùng chịu hoặc số tiền, phần chia được tính lại toàn bộ

## Ví dụ đầy đủ

Nhóm 4 người: An, Bình, Chi, Dũng. Bữa ăn 1.000.000.

- An và Bình trả hộ (hai người trả)
- Cả bốn cùng chịu
- Dũng chịu cố định 100.000 (ăn ít)
- Ba người còn lại chia 900.000 → mỗi người 300.000

Kết quả phải chịu: An 300.000, Bình 300.000, Chi 300.000, Dũng 100.000.  
Phần đã trả: An 500.000, Bình 500.000, Chi 0, Dũng 0.

Số dư:

- An: +200.000 (được nhận)
- Bình: +200.000 (được nhận)
- Chi: −300.000 (đang nợ)
- Dũng: −100.000 (đang nợ)

Ứng dụng sẽ gợi ý Chi và Dũng chuyển tiền cho An và Bình để cân bằng.
