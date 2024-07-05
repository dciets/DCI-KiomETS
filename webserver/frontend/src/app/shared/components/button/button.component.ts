import { Component, Input, OnInit } from '@angular/core';

@Component({
  selector: 'ui-button',
  templateUrl: './button.component.html',
  styleUrls: ['./button.component.css'],
})
export class ButtonComponent implements OnInit {
  @Input() onClick!: () => void;
  @Input() label!: string;

  constructor() {}

  ngOnInit(): void {}
}
