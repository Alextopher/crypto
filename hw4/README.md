# HW 4

This is elliptic curve elgamal. We were given all of the private and public key and we're just asked to decrypt a message. Every line in [a4.cipher](a4.cipher) is in the form

```
CipherX CipherY HalfmaskX HalfmaskY
CipherX CipherY HalfmaskX HalfmaskY
CipherX CipherY HalfmaskX HalfmaskY
CipherX CipherY HalfmaskX HalfmaskY
CipherX CipherY HalfmaskX HalfmaskY
...
```

Each message encodes a single character.

## Result

> a = 31969280088361358982891392020418160677752584268704094389910371947753244946933347935878625023919402452740316172441746808587155134252287220284210084660983594234250923217292376810981985895141552247204162868475624471518794404318233925012571130473382462504642943159579216921221327811371889912608373396443367093840

This is the private key for [hw3](../hw3)