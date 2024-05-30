# 5 Tühistus- ja ennistusnimekiri

VIS3 edastab EHS-le tühistus- ja ennistusnimekirja (lühidalt - tühistusnimekiri).

Tühistusnimekiri sisaldab andmeid isikute kohta, kelle e-hääl tuleb tühistada (on paberhääletanud ja sellest tulenevalt ei lähe e-hääl arvesse valimistulemuste kokkulugemisel) või ennistada (s.t. tühistatakse eelnev tühistamine ning häälte uuesti üle lugemisel võetakse ennistatud e-hääl arvesse).

Fail on JSON-formaadis.

## Protseduur

Tühistusnimekiri edastakse inim-masin-protseduuriga:

1) EHS operaator laeb JSON-faili EHS-st alla;

2) EHS operaator allkirjastab faili digitaalselt väljaspool EHS-i ja annab selle VIS peakasutajale;

3)  VIS peakasutaja laeb digikonteinerist välja võetud faili VIS3-i.

## Edastatav fail

Aluseks on EHS senine liides VIS2-ga, mis on spetsifitseeritud dokumendis "IVXV protokollide kirjeldus" (v 1.5.0, IVXV-PR-1.5.0, 20.04.2019), [https://www.valimised.ee/sites/default/files/uploads/eh/IVXV-protokollid.pdf](https://www.valimised.ee/sites/default/files/uploads/eh/IVXV-protokollid.pdf), jaotises "Tühistus- ja ennistusnimekiri" (jaotis 9.1).

Valimissündmuse identifikaator peab vastama formaadile [Valimissündmuse identifikaator](../valimissündmuse_identifikaator.md).

Faili struktuur (JSON-skeem): [revoke.schema](revoke.schema)

Faili näited (JSON):

- [revoke_EP.json](revoke_EP.json)
- [revoke_KOV.json](revoke_KOV.json)
- [revoke_RH.json](revoke_RH.json)
- [revoke_RK.json](revoke_RK.json)
