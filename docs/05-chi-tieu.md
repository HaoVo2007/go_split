---
project: go-split
topic: chi-tieu
language: vi
---

# Khoản chi

Khoản chi là một hóa đơn chung trong nhóm: bữa ăn, vé xe, tiền phòng...

## Thông tin cần có khi thêm khoản chi

- Thuộc nhóm nào
- Tên khoản chi
- Số tiền tổng
- Ngày (ví dụ 2026-08-19)
- Loại chi tiêu (chuỗi tự đặt: ăn uống, đi lại, nhà cửa...)
- Người trả tiền — một hoặc nhiều người
- Người cùng chịu — danh sách người chia khoản này
- Ảnh hóa đơn (không bắt buộc)
- Cách chia tùy chỉnh (không bắt buộc)

**Quy tắc bắt buộc:** mọi người trả tiền phải nằm trong danh sách người cùng chịu. Ví dụ An trả tiền thì An cũng phải là người cùng chịu khoản đó.

## Người trả và người cùng chịu khác nhau thế nào?

- **Người trả** là người đưa tiền cho cửa hàng / ứng trước.
- **Người cùng chịu** là những người thật sự hưởng khoản chi, sẽ gánh một phần số tiền.

Một người có thể vừa trả vừa cùng chịu. Nhiều người có thể cùng trả một hóa đơn.

## Xem khoản chi

Khi xem một khoản chi, bạn thấy tên, số tiền, loại, ảnh, danh sách người trả, người cùng chịu, và phần chia tùy chỉnh (nếu có), kèm thời điểm tạo và cập nhật.

Danh sách khoản chi theo nhóm được phân trang. Mặc định mỗi trang 3 khoản, trang đầu là trang 1. Khoản mới hơn hiện trước.

## Sửa khoản chi

Có thể đổi ngày, tên, số tiền, loại, ảnh, người trả, người cùng chịu, và cách chia.

Khi đổi danh sách người cùng chịu hoặc cách chia, hệ thống tính lại phần của từng người từ đầu.

Khi đổi ảnh, ảnh cũ được thay bằng ảnh mới.

## Xóa khoản chi

Xóa là xóa mềm: khoản chi không còn hiện trong danh sách và không còn tính vào số dư. Phần chia gắn với khoản đó cũng bị gỡ.

## Ảnh hóa đơn

Ảnh giúp mọi người nhớ hóa đơn. Không bắt buộc. Ảnh được lưu trên dịch vụ ảnh của ứng dụng.
