import { Component, OnInit } from '@angular/core';

@Component({
  selector: 'app-home',
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.css'],
})
export class HomeComponent implements OnInit {
  language: 'fr'|'en';
  constructor() {
    this.language = (localStorage.getItem("language") ?? "fr") as 'fr' | 'en';
  }

  changeLanguage(language: 'fr'|'en') {
    this.language = language;
    localStorage.setItem("language", language);
  }

  ngOnInit(): void {}
}
