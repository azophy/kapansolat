KapanSolat
==========

straightforward website to get prayer time (shalat) info

## features
- [ ] display TUI respond when called with curl
- [ ] display json formatted response when requested
- [ ]
- [ ] .
- [ ] .
- [ ] .
- [ ] .

no-nonsense jadwal solat. TUI. jadi bisa


curl kapansolat.com

return:
perkiraan lokasi sekarang
jam skrg

count down ke waktu solat terdekat

jadwal solat hari ini

logic routing:
- kalau ada slug cari dulu
-   found, tampilkan
-   not found: next
- kalau ada cookie, tampilkan dulu
- gak ada data, search from scratch
-   geocode latlong by ip
-   store result as cookie

baru display waktu
