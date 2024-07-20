import { Component, OnInit } from '@angular/core';
import {environment} from "../../environments/environment";

interface Score {
  rank: number;
  name: string;
  score: number;
}

@Component({
  selector: 'app-scoreboard',
  templateUrl: './scoreboard.component.html',
  styleUrls: ['./scoreboard.component.css'],
})
export class ScoreboardComponent implements OnInit {
  scores?: Score[];
  private socket: WebSocket|undefined;

  constructor() {}

  ngOnInit(): void {
    this.socket = new WebSocket('ws://'+environment.serverAddr+'/ws/scoreboard');
    this.socket.onmessage = (ev: MessageEvent) => {const decoded = atob(ev.data.trim());
      if (decoded.trim() !== '') {
        this.scores = (JSON.parse(decoded) as any[])
          .sort((a, b) => b.score - a.score).map((data, index) => ({...data, rank: index + 1}));
      }
    }
  }
}
