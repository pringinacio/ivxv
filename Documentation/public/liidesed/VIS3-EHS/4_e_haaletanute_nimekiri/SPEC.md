# E-hääletanute nimekiri

## Protseduur

E-hääletanute nimekiri on pärast e-hääletamise lõppu EHS poolt genereeritav ja väljastatav nimekiri e-hääletanud isikutest. E-hääletanute nimekiri loetakse sisse VIS3-e (moodul NIM).

E-hääletanute nimekiri edastakse inim-masin-protseduuriga:

1) EHS operaator laeb JSON-faili EHS-st alla;

2) EHS operaator allkirjastab faili digitaalselt väljaspool EHS-i ja annab selle VIS peakasutajale;

3)  VIS peakasutaja laeb digikonteinerist välja võetud faili VIS3-i.

Fail on JSON-formaadis. Faili näidised on tööülesandes.

Faili struktuur on kirjeldatud JSON-skeemiga.

## Edastatav fail

Aluseks on EHS senine liides VIS2-ga, mis on spetsifitseeritud dokumendis "IVXV protokollide kirjeldus" (v 1.5.0, IVXV-PR-1.5.0, 20.04.2019), [https://www.valimised.ee/sites/default/files/uploads/eh/IVXV-protokollid.pdf](https://www.valimised.ee/sites/default/files/uploads/eh/IVXV-protokollid.pdf), jaotises "E-hääletanute nimekiri" (jaotis 9.2).

Faili struktuur (JSON-skeem): [onlinevoters.schema](onlinevoters.schema)

Valimissündmuse identifikaator peab vastama formaadile [Valimissündmuse identifikaator](../valimissündmuse_identifikaator.md).

Faili näited (JSON):

- [onlinevoters_EP.json](onlinevoters_EP.json)
- [onlinevoters_KOV.json](onlinevoters_KOV.json)
- [onlinevoters_RH.json](onlinevoters_RH.json)
- [onlinevoters_RK.json](onlinevoters_RK.json)
