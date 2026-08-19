---
project: go-split
topic: so-du-quyet-toan
language: vi
---

# Số dư và quyết toán

go-split không chuyển tiền thật. Nó chỉ nói **ai nên trả cho ai bao nhiêu**. Mọi người tự chuyển ngoài ứng dụng rồi coi như đã tất toán.

## Số dư trong một khoản chi

Xem quyết toán theo từng hóa đơn khi chỉ muốn xử lý một bữa / một vé.

Với mỗi người liên quan, ứng dụng cho biết:

- Đã trả bao nhiêu
- Phải chịu bao nhiêu
- Số dư
- Trạng thái: đang được nhận (creditor), đang nợ (debtor), hoặc đã cân (settled)

Kèm danh sách gợi ý: người nợ chuyển cho người được nhận, số tiền cụ thể.

## Số dư cả nhóm

Xem số dư nhóm khi muốn tất toán cả tháng / cả chuyến đi.

Gồm:

- Tổng mọi khoản chi chưa xóa trong nhóm
- Từng thành viên (kể cả chủ nhóm): đã trả, phải chịu, số dư, trạng thái
- Phần **bạn đang xem** đã trả và phải chịu
- Danh sách gợi ý chuyển tiền để cả nhóm về cân bằng

Chỉ tính các khoản chi còn hiệu lực. Khoản đã xóa không còn trong số dư.

## Cách hiểu gợi ý chuyển tiền

Gợi ý theo hướng: **người đang nợ** trả **cho người đang được nhận**.

Mục tiêu là ít giao dịch nhất có thể, không phải mỗi người trả từng hóa đơn cũ.

Ví dụ An +200, Bình +100, Chi −300. Gợi ý có thể là Chi chuyển 200 cho An và 100 cho Bình. Không cần Chi trả từng bữa ăn trong quá khứ.

## Khi nào số dư đổi?

Số dư đổi khi:

- Thêm khoản chi mới
- Sửa số tiền, người trả, người cùng chịu, cách chia
- Xóa khoản chi

Không tự đổi vì thời gian trôi. Không có nút “đã trả” trong hệ thống: sau khi chuyển tiền ngoài đời, nhóm có thể xóa / ngừng dùng, hoặc tự theo dõi đã trả. Ứng dụng không ghi nhận giao dịch chuyển khoản ngân hàng.

## Chủ nhóm và thành viên

Số dư nhóm tính cho chủ nhóm và mọi thành viên đã mời, kể cả người chưa từng có tên trên một hóa đơn (họ sẽ 0 nếu chưa tham gia khoản nào).
