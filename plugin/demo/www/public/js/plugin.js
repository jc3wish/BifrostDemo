// JavaScript Document
var Plugin = ( function () {

    function sayHello(){
        console.log( this ); /* this 指向 window */
        console.log( '欢迎访问nDos的博客' );
    }
    return class Class01{
        constructor( v = '初始值' ){
            pVal = v;
        }
        get val(){
            sayHello();
            return pVal;
        }
        set val( v ){
            pVal = v;
        }
    };
} )();