# 2 Valikute nimekiri (kandidaatide nimekiri)

## Protseduur

VIS3 edastab EHS-le andmed registreeritud kandidaatide (valimistel) või vastusevariantide (rahvahääletusel) kohta.

Andmed edastatakse inim-masin-protseduuriga:

1. VIS peakasutaja laeb JSON faili VIS3-st alla; Kasutaja peab valima aktiivsete valimissündmuste seast soovitud sündmuse, mille kohta infot alla laadida. Valimissündmus on aktiivne, kui tema staatus pole `closed` ega `deleted`.
2. VIS peakasutaja allkirjastab faili digitaalselt väljaspool VIS-i ja annab selle EHS operaatorile;
3. EHS operaator laeb digiallkirjastatud faili EHS-i.

Fail on JSON-formaadis (näidis lisatud tööülesandele). Faili ei allkirjastata.

## Edastatav fail

Aluseks on EHS senine liides VIS2-ga, mis on spetsifitseeritud dokumendis "IVXV protokollide kirjeldus" (v 1.5.0, IVXV-PR-1.5.0, 20.04.2019), [https://www.valimised.ee/sites/default/files/uploads/eh/IVXV-protokollid.pdf](https://www.valimised.ee/sites/default/files/uploads/eh/IVXV-protokollid.pdf), jaotises "Valikute nimekiri" (jaotis 3.4).

Valimissündmuse identifikaator peab vastama formaadile [Valimissündmuse identifikaator](../valimissündmuse_identifikaator.md).

Faili struktuur (JSON-skeem): [choices.schema](choices.schema)

Faili näited (JSON):

- [choices_KOV.json](choices_KOV.json)
- [choices_EP.json](choices_EP.json)
- [choices_RK.json](choices_RK.json)

## Taustainfo

Valikute nimekiri sisaldab andmeid kandidaatide (valimistel) või vastusevariantide (rahvahääletusel)
kohta. Valimiste korral on lisaks kandidaadi andmetele nimekirjas ka tema erakonna või valimisliidu nimi, mille nimekirjas ta kandideerib.

Valijale elektroonilise hääletamise käigus nähtavaid valimiste vahelisi süsteemseid erinevusi on kaks:

1. Rahvahääletusel ei valita erakondadesse kuuluvate kandidaatide vahel vaid vastatakse „JAH“/“EI“ rahvahääletuse küsimusele;

2. Riigikogu, KOV ja Euroopa Parlamendi valimistel antakse hääl ühele kandidaadile, kes võib, aga ei pruugi kuuluda poliitilise ühenduse nimekirja.

Protokollistik kodeerib valija võimalikud valikud ringkonnas kuni 11-kohalise arvväärtusena, mis valikute nimekirjas kodeeritakse koos ringkonna EHAK-koodiga. Valijale tohivad kättesaadavad olla ainult tema ringkonnakohased valikud. Valijarakendus peab seda omadust tagama ning hääletamistulemust arvutav rakendus kontrollima.
