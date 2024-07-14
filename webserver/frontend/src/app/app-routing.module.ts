import { NgModule } from '@angular/core';
import { HomeComponent } from './home/home.component';
import { ScoreboardComponent } from './scoreboard/scoreboard.component';
import { RouterModule, Routes } from '@angular/router';
import { SignupComponent } from './signup/signup.component';
import { GameComponent } from './game/game.component';
import {AdminComponent} from "./admin/admin.component";
import {FormsModule} from "@angular/forms";

const routes: Routes = [
  {
    path: '',
    component: HomeComponent,
  },
  {
    path: 'signup',
    component: SignupComponent,
  },
  {
    path: 'game',
    component: GameComponent,
  },
  {
    path: 'scoreboard',
    component: ScoreboardComponent,
  },
  {
     path: 'Admin',
     component: AdminComponent,
  },
  { path: '**', redirectTo: 'home', pathMatch: 'full' }, // if route doesn't exist
];

@NgModule({
  declarations: [],
  imports: [
    RouterModule.forRoot(routes),
    FormsModule
  ],
  exports: [RouterModule],
})
export class AppRoutingModule {}
