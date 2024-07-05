import { Component, OnInit } from '@angular/core';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-game',
  templateUrl: './game.component.html',
  styleUrls: ['./game.component.css'],
})
export class GameComponent implements OnInit {
  private socket: WebSocket|undefined

  constructor() {}

  ngOnInit(): void {
    this.socket = new WebSocket('ws://'+environment.serverAddr+'/ws/game');
    this.socket.onmessage = (ev: MessageEvent) => {
      console.log(ev.data);
    }
  }
}
