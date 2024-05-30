#!/bin/bash

TOTAL=$(($1-1))

source card_options.conf

function erase {
    # erase-card will explicitly ignore any PIN from commandline/env, must always enter manually
    echo $PKCS15_SO_PIN | pkcs15-init --erase-card -r $1 ;
}

function create {
    pkcs15-init --pin "env:PKCS15_PIN" --puk "env:PKCS15_PUK" --so-pin "env:PKCS15_SO_PIN" --so-puk "env:PKCS15_SO_PUK" --create-pkcs15 -r $1;
    pkcs15-init --pin "env:PKCS15_PIN" --puk "env:PKCS15_PUK" --auth-id 01 --store-pin --label el -r $1;
    pkcs15-init --finalize -r $1;
}

echo "PIN will be $PKCS15_PIN."
echo "During initalization, the pkcs15-init script asks for Security Officer PIN."
echo "This is provided automagically, you don't need to insert anything".
echo "Please wait until erasing and initalizing $1 smart cards"
echo "==============================="

for i in `seq 0 $TOTAL`; do
  echo "Erasing card $i"
    erase $i;
  echo "Creating card $i"
    create $i;
  echo "Done"
  echo
done

echo "==============================="
echo "All done."
echo "PIN is $PKCS15_PIN."
