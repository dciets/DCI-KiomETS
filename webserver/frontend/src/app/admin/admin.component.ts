import {Component, OnInit} from '@angular/core';
import {FormsModule} from "@angular/forms";
import {CommonModule} from "@angular/common";
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-admin',
  templateUrl: './admin.component.html',
  styleUrls: ['./admin.component.css'],
  standalone: true,
  imports: [FormsModule, CommonModule]
})
export class AdminComponent implements OnInit {
  parameters: any = {
    mapSize: 0,
    soldierSpeed: 0,
    soldierCreationSpeed: 3,
    terrainChangeSpeed: 0,
    gameLength: 0,
  }

  backendResponse: string = "";

  players:{Name:string, UID:string}[]= []
  ngOnInit(): void {
    this.GetParameters(null);
    this.GetPlayers(null);
  }

  GetParameters(event: Event | null){
    if (event){
      this.flashElement(event.target as HTMLElement)
    }
    // fetch parameters from server
    fetch(environment.serverAddr).then((response) => {
      response.json().then((data) => {
        console.log(data);
        this.parameters = data;
      })
    })
  }

  GetPlayers(event: Event | null){
    if (event){
      this.flashElement(event.target as HTMLElement)
    }
    // fetch players from server
    fetch(environment.serverAddr).then((response) => {
      response.json().then((data) => {
        console.log(data);
        this.players = data;
      })
    })
  }

  SetParameters(event: Event){
    this.flashElement(event.target as HTMLElement)
    // send parameters to server
    fetch(environment.serverAddr, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(this.parameters),
    }).then((response) => {
      response.json().then((data) => {
        console.log(data);
        this.parameters = data;
      })
    })
  }

  flashElement(element: HTMLElement){
    element.style.boxShadow = '0 0 10px 5px #fff2';
    setTimeout(() => {
      element.style.boxShadow = 'none';
    }, 300);
  }
  startGame(){
    // start game
    fetch(environment.serverAddr, {
      method: 'POST',
    }).then((response) => {
      response.text().then((data) => {
        console.log(data);
        this.backendResponse = `${response.status} ${data}`
      })
    })
  }
  endGame(){
    // end game
    fetch(environment.serverAddr, {
      method: 'POST',
    }).then((response) => {
      response.text().then((data) => {
        console.log(data);
        this.backendResponse = `${response.status} ${data}`
      })
    })
  }
  status(){
    // get game status
    fetch(environment.serverAddr).then((response) => {
      response.text().then((data) => {
        console.log(data);
        this.backendResponse = `${response.status} ${data}`
      })
    })
  }


}
