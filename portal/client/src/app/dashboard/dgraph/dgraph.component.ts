
import { Component, HostListener, OnInit, OnDestroy, Input, ElementRef, ViewChild} from '@angular/core';
import { Graph } from '../../models/graph.model';
import { DashboardService } from '../services/dashboard.service'
import { MenuService } from '../../services/menu.service'
import * as d3 from 'd3';
import { GraphText } from '../services/graph-text.service'
import { GraphPie } from '../services/graph-pie.service'
import { GraphLines } from '../services/graph-lines.service'
import { GraphBars} from '../services/graph-bars.service'
import { GraphAreas } from '../services/graph-areas.service'
import { GraphBubbles } from '../services/graph-bubbles.service'

@Component({
  selector: 'app-dgraph',
  templateUrl: "./dgraph.component.html",
  styleUrls: ['./dgraph.component.css'],
  providers: [ GraphText, GraphPie, GraphLines, GraphBars, GraphAreas, GraphBubbles ]


})
export class DGraphComponent implements OnInit, OnDestroy {
  @ViewChild('chart') private chartContainer: ElementRef;
  @Input() public graph : Graph;
  //Resizer variable
  px: number = 0;
  py: number = 0;
  minWidth: number = 20;
  minHeight: number = 20;
  draggingCorner: boolean = false;
  draggingWindow: boolean = false;
  resizer: Function;
  private serviceMap = {}

  constructor(
    public dashboardService : DashboardService,
    private menuService : MenuService,
    private graphPie : GraphPie,
    private graphLines : GraphLines,
    private graphBars : GraphBars,
    private graphAreas : GraphAreas,
    private graphBubbles : GraphBubbles,
    private graphText : GraphText) {
    this.serviceMap['text'] = graphText;
    this.serviceMap['pie'] = graphPie;
    this.serviceMap['lines'] = graphLines;
    this.serviceMap['bars'] = graphBars;
    this.serviceMap['areas'] = graphAreas;
    this.serviceMap['bubbles'] = graphBubbles;
  }

  ngOnInit() {
    this.serviceMap[this.graph.type].init(this.graph, this.chartContainer)
    this.dashboardService.onNewData.subscribe(
      () => {
        this.updateGraph();
      }
    )
    this.menuService.onWindowResize.subscribe(
      (win) => {
        this.resizeGraph()
      }
    );
  }

  ngOnDestroy() {
    this.serviceMap[this.graph.type].destroy()
    //this.metricsService.onNewData.unsubscribe()
  }

  createGraph() {
    this.serviceMap[this.graph.type].createGraph(this.graph, this.chartContainer)
  }

  clearGraph() {
    this.serviceMap[this.graph.type].clearGraph()
  }

  resizeGraph() {
    this.serviceMap[this.graph.type].resizeGraph(this.graph, this.chartContainer)
  }

  updateGraph() {
    this.serviceMap[this.graph.type].updateGraph(this.graph)
  }


//----------------------------------------------------------------------------
// Edition mode
//----------------------------------------------------------------------------

  onWindowPress(event: MouseEvent) {
    this.draggingWindow = true;
    this.px = event.clientX;
    this.py = event.clientY;
    this.dashboardService.selected = this.graph
    //event.stopPropagation();
    return false
  }

  onWindowDrag(event: MouseEvent) {
    if (!this.draggingWindow) {
      return;
    }
    let offsetX = event.clientX - this.px;
    let offsetY = event.clientY - this.py;
    this.graph.x += offsetX;
    this.graph.y += offsetY;
    this.px = event.clientX;
    this.py = event.clientY;
    event.stopPropagation();
  }

  topLeftResize(offsetX: number, offsetY: number) {
    this.graph.x += offsetX;
    this.graph.y += offsetY;
    this.graph.width -= offsetX;
    this.graph.height -= offsetY;
    this.resizeGraph()
  }

  topRightResize(offsetX: number, offsetY: number) {
    this.graph.y += offsetY;
    this.graph.width += offsetX;
    this.graph.height -= offsetY;
    this.resizeGraph()
  }

  bottomLeftResize(offsetX: number, offsetY: number) {
    this.graph.x += offsetX;
    this.graph.width -= offsetX;
    this.graph.height += offsetY;
    this.resizeGraph()
  }

  bottomRightResize(offsetX: number, offsetY: number) {
    this.graph.width += offsetX;
    this.graph.height += offsetY;
    this.resizeGraph()
  }

  onCornerClick(event: MouseEvent, resizer?: Function) {
    this.draggingCorner = true;
    this.dashboardService.selected = this.graph
    this.px = event.clientX;
    this.py = event.clientY;
    this.resizer = resizer;
    event.preventDefault();
    event.stopPropagation();
  }

  @HostListener('document:mousemove', ['$event'])
    onCornerMove(event: MouseEvent) {
      if (!this.draggingCorner) {
        return;
      }
      let offsetX = event.clientX - this.px;
      let offsetY = event.clientY - this.py;

      let lastX = this.graph.x;
      let lastY = this.graph.y;
      let pWidth = this.graph.width;
      let pHeight = this.graph.height;

      this.resizer(offsetX, offsetY);
      if (this.graph.width < this.minWidth || this.graph.height < this.minHeight) {
        this.graph.x = lastX;
        this.graph.y = lastY;
        this.graph.width = pWidth;
        this.graph.height = pHeight;
        this.resizeGraph()
      }
      this.px = event.clientX;
      this.py = event.clientY;
    }

  @HostListener('document:mouseup', ['$event'])
    onCornerRelease(event: MouseEvent) {
      this.draggingWindow = false;
      this.draggingCorner = false;
      //event.stopPropagation();
      return false
    }

}
