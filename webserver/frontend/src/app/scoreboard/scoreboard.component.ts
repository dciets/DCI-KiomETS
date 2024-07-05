import { Component, OnInit } from '@angular/core';

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

  constructor() {}

  ngOnInit(): void {
    // TODO get scores from backend
    this.scores = [
      { rank: 1, name: 'Team 1', score: 100 },
      { rank: 2, name: 'Team 2', score: 90 },
      { rank: 3, name: 'Team 3', score: 80 },
      { rank: 4, name: 'Team 4', score: 70 },
      { rank: 5, name: 'Team 5', score: 60 },
      { rank: 6, name: 'Team 6', score: 50 },
      { rank: 7, name: 'Team 7', score: 40 },
      { rank: 8, name: 'Team 8', score: 30 },
      { rank: 9, name: 'Team 9', score: 20 },
      { rank: 10, name: 'Team 10', score: 10 },
    ];
  }
}
