# Bank
This is a moc bank server application on golang
to build and run in just run 'make run' in project directory

server runs on localhost:3000

RestApi:

//get all acc
GET localhost:3000/account

//need to use in header 'x-jwt-token'  with of acc
//{
//"number": 582384,
//"password": "password"
//}
GET localhost:3000/login

//create account use in body
//example :

//{
//"firstName": "Alex",
//"lastName": "Tek",
//"password": "ultrapassword"
//}
POST localhost:3000/account

//delete account
DELETE localhost:3000/account/{id}
