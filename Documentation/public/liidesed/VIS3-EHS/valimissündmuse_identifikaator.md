# Valimissündmuse identifikaator

Üht valimist puudutav andmestik on seotud unikaalse valimissündmuse identifikaatori abil.

Valimissündmuse identifikaator on sõne, mis koosneb kahest kohustuslikust ja kahest valikulisest osast:
1. valimissündmuse tüüp, vastavalt allolevale tabelile, nt `RK`.
2. aastanumbrist (nt `2023`)
3. valikulisest erakorraliste valimiste tunnusest `_E` ja
4. valikulisest järjenumbrist <nr>, nt `2`.

Ülalnimetatud osad eraldatakse üksteisest allkriipsudega `_`.

Järjenummerdatakse tüüpide kaupa. Seejuures aasta esimese samatüübilise valimissündmuse korral järjenumbrit ei näidata.

Näited:
- `KOV_2021`
- `RH_2021`
- `RH_2021_2` (2021. a teine rahvahääletus)
- `RK_2023`
- `RK_2023_E` (2023. a Riigikogu erakorralised valimised)
- `RK_2023_E_2` (2023. a Riigikogu teised erakorralised valimised).

EHSi jaoks on valimise identifikaator kuni 28 ASCII-tähemärgi pikkune sõne ning eelnevalt kirjeldatud struktuuri põhjal EHS otsuseid ei tee.

## Valimissündmuse tüüp

Valimissündmuse tüüp esitatakse koodiga:

- `KOV` - Kohaliku omavalitsuse volikogu valimised
- `EP` - Euroopa Parlamendi valimised
- `RK` - Riigikogu valimised
- `RH` - Rahvahääletus.

