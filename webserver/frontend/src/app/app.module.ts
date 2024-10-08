import { NgModule } from '@angular/core';
import { BrowserModule } from '@angular/platform-browser';

import { AppComponent } from './app.component';
import { AppRoutingModule } from './app-routing.module';
import { HomeComponent } from './home/home.component';
import { ScoreboardComponent } from './scoreboard/scoreboard.component';
import { ButtonComponent } from './shared/components/button/button.component';
import { SignupComponent } from './signup/signup.component';
import { GameComponent } from './game/game.component';
import {CommonModule, NgOptimizedImage} from "@angular/common";
import {ReactiveFormsModule,FormsModule} from "@angular/forms";

@NgModule({
  declarations: [
    AppComponent,
    HomeComponent,
    ScoreboardComponent,
    ButtonComponent,
    SignupComponent,
    GameComponent,
  ],
    imports: [BrowserModule, AppRoutingModule, CommonModule, ReactiveFormsModule, FormsModule, NgOptimizedImage],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
