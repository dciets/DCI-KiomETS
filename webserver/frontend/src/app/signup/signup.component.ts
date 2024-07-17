import { Component, OnInit } from '@angular/core';
import {environment} from "../../environments/environment";

@Component({
  selector: 'app-signup',
  templateUrl: './signup.component.html',
  styleUrls: ['./signup.component.css']
})
export class SignupComponent implements OnInit {

  botName: string = "";
  uid: string = "";
  error: string = "";
  copied: string = "";

  constructor() { }

  ngOnInit(): void {
  }

  signup(): void {
    fetch("http://"+environment.serverAddr  +"/api/agent", {
      method: 'POST',
      body: JSON.stringify({name: this.botName}),
      headers: {
        'Content-Type': 'application/json'
      }
    }).then((response) => {
      if (response.status === 200) {
        response.json().then((data) => {
          this.uid = data.UID;
        })
      } else {
        response.text().then((data) => {
          this.error = data;
        })
      }
    })
  }

  copyUID(): void {
    navigator.clipboard.writeText(this.uid).then(() => {
      this.copied = "Copied!";
      setTimeout(() => {
        this.copied = "";
      }, 1000);
    }, (err) => {
      console.error('Async: Could not copy text: ', err);
    });
  }

}
