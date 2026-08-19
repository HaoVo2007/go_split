---
project: go-split
topic: bang-dieu-khien
language: vi
---

# Bảng điều khiển cá nhân

Bảng điều khiển (dashboard) là màn hình tổng hợp **của bạn**, gom mọi nhóm bạn tạo hoặc được mời.

## Bạn thấy gì?

### Số dư cá nhân

- Phần bạn phải chịu (you owed / you_owed)
- Phần bạn đã trả / đã ứng (you paid / you_paid)
- Chênh lệch: đã trả trừ đi phần phải chịu

Dương nghĩa overall bạn đang cho người khác nợ. Âm nghĩa overall bạn đang nợ.

### Tổng quan

- Số nhóm bạn đang tham gia
- Số khoản chi (giao dịch) trong các nhóm đó
- Số bạn bè: những người khác xuất hiện trong các nhóm của bạn, không tính chính bạn

### Chi tiêu của bạn

- Tổng đã trả
- Tổng phần được chia cho bạn (total shared)

Hai số này cùng nguồn với số dư ở trên.

### Thống kê nổi bật

- **Nhóm chi nhiều nhất:** nhóm có tổng tiền các khoản chi lớn nhất trong các nhóm của bạn
- **Bạn đồng hành nhiều nhất:** người (không phải bạn) xuất hiện trong nhiều khoản chi nhất với bạn

Nếu chưa có nhóm, mọi số là 0. Nếu có nhóm nhưng chưa có khoản chi, số nhóm vẫn hiện, còn giao dịch và thống kê chi tiêu là 0.

## Khi nào dùng dashboard, khi nào vào nhóm?

- Dashboard: “tôi đang nợ hay đang được nhận trên toàn bộ app?”
- Số dư nhóm: “chuyến Đà Lạt này ai nợ ai?”
- Quyết toán một khoản: “bữa tối hôm qua chia xong chưa?”

## Lưu ý

Dashboard chỉ tính nhóm bạn còn thấy (chưa xóa) và khoản chi chưa xóa. Người bạn trong thống kê lấy tên hồ sơ; chưa có tên thì hiện email.
