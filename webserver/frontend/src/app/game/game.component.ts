import {AfterViewInit, Component, ElementRef, OnInit, ViewChild} from '@angular/core';
import { environment } from 'src/environments/environment';

@Component({
  selector: 'app-game',
  templateUrl: './game.component.html',
  styleUrls: ['./game.component.css'],
})
export class GameComponent implements OnInit, AfterViewInit {
  private socket: WebSocket|undefined;
  private distanceMultiplier: number = 30;
  private displacement: {x: number, y: number} = {x: 0, y: 0};
  protected data: any;
  private isDragging: boolean = false;

  @ViewChild('canvas') canvas: ElementRef|undefined;

  constructor() {}

  draw() {
    const canvas = (this.canvas?.nativeElement as HTMLCanvasElement);
    const context: CanvasRenderingContext2D = canvas.getContext('2d')!;
    const backgroundColor = '#000';
    const neutralColor = '#fff';

    const centerX = canvas.width / 2 + this.displacement.x;
    const centerY = canvas.height / 2 + this.displacement.y;

    const sidesCoordinates = [[1, 0], [0.5, 0.866], [-0.5, 0.866], [-1, 0], [-0.5, -0.866], [0.5, -0.866]];

    context.fillStyle = backgroundColor;
    context.fillRect(0, 0, canvas.width, canvas.height);

    context.lineWidth = 1;
    context.font = '16px arial';

    for (const pipe of this.data.pipes) {
      const terrain1Positions = this.data.terrains[pipe.first].position;
      const terrain2Positions = this.data.terrains[pipe.second].position;
      const terrain1X = terrain1Positions[0] * this.distanceMultiplier + centerX;
      const terrain1Y = terrain1Positions[1] * this.distanceMultiplier + centerY;
      const terrain2X = terrain2Positions[0] * this.distanceMultiplier + centerX;
      const terrain2Y = terrain2Positions[1] * this.distanceMultiplier + centerY;

      context.strokeStyle = neutralColor;
      context.beginPath();
      context.moveTo(terrain1X, terrain1Y);
      context.lineTo(terrain2X, terrain2Y);
      context.stroke();

      for (const soldier of pipe.soldiers) {
        const length = soldier.length * 1.0 / pipe.length;
        const x = (terrain2X - terrain1X) * length + terrain1X;
        const y = (terrain2Y - terrain1Y) * length + terrain1Y;

        context.strokeStyle = this.data.players[soldier.ownerIndex].color;
        context.fillStyle = context.strokeStyle;
        context.fillRect(x - 1, y - 1, 2, 2);

        if (this.distanceMultiplier > 12) {
          context.fillText(soldier.soldierCount, x + this.distanceMultiplier, y);
        }
      }
    }

    context.fillStyle = neutralColor;

    for (const terrain of this.data.terrains) {
      if (terrain.ownerIndex !== -1) {
        context.fillStyle = this.data.players[terrain.ownerIndex].color;
      } else {
        context.fillStyle = neutralColor;
      }

      const terrainX = terrain.position[0] * this.distanceMultiplier + centerX;
      const terrainY = terrain.position[1] * this.distanceMultiplier + centerY;
      context.beginPath();
      context.moveTo(
        terrainX + sidesCoordinates[0][0] * Math.max(20, this.distanceMultiplier) / 3,
        terrainY + sidesCoordinates[0][1] * Math.max(20, this.distanceMultiplier) / 3);

      for (const side of sidesCoordinates) {
        context.lineTo(
          terrainX + side[0] * Math.max(20, this.distanceMultiplier) / 3,
          terrainY + side[1] * Math.max(20, this.distanceMultiplier) / 3);
      }
      context.fill();

      if (this.distanceMultiplier > 12) {
        context.fillText(terrain.numberOfSoldier.toString(), terrainX + this.distanceMultiplier, terrainY);
      }
    }
  }

  scroll(amount: number) {
    this.distanceMultiplier = Math.max(0, Math.min(100, this.distanceMultiplier + amount));
    this.draw();
    localStorage.setItem('zoom', this.distanceMultiplier.toString());
  }

  startDrag() {
    this.isDragging = true;
  }

  stopDrag() {
    this.isDragging = false;
    localStorage.setItem('displacement', JSON.stringify(this.displacement));
  }

  drag(event: MouseEvent) {
    if (this.isDragging) {
      this.displacement.x += event.movementX;
      this.displacement.y += event.movementY;
      this.draw();
    }
  }

  ngOnInit(): void {
    if (localStorage.getItem('zoom')) {
      this.distanceMultiplier = parseInt(localStorage.getItem('zoom')!)
    }
    if (localStorage.getItem('displacement')) {
      this.displacement = JSON.parse(localStorage.getItem('displacement')!);
    }
  }

  ngAfterViewInit(): void {
    const canvas = (this.canvas?.nativeElement as HTMLCanvasElement);
    canvas.height = canvas.getBoundingClientRect().height / canvas.getBoundingClientRect().width * canvas.width;

    this.socket = new WebSocket('ws://'+environment.serverAddr+'/ws/game');
    this.socket.onmessage = (ev: MessageEvent) => {
      if (ev.data.trim() !== '') {
        this.data = JSON.parse(ev.data);
        this.draw();
      }
    }
  }
}
