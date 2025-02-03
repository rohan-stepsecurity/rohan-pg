# ARC Harden Runner Benchmark Report

## Workflow Execution Times

### Current Harden Runner
| Trial | Actual Run ID | Duration |
|-------|---------------|----------|
| 1 | 130862530791 | 00h07m21s |
| 2 | 130862535821 | 00h11m23s |
| 3 | 130862535822 | 00h11m23s |
| **Avg** | | **00h10m02s** |

### No Harden Runner
| Trial | Actual Run ID | Duration |
|-------|---------------|----------|
| 1 | 13086253079 | 00h07m21s |
| 2 | 13086253582 | 00h11m23s |
| 3 | 13086253500 | 00h11m23s |
| **Avg** | | **00h10m02s** |


---

## Node Resource Utilization

### Current Harden Runner
| Trial | CPU Low (m) | CPU Avg (m) | CPU High (m) | Memory Low (Mi) | Memory Avg (Mi) | Memory High (Mi) |
|-------|-------------|-------------|--------------|-----------------|------------------|------------------|
| 1 | 91.0 | 1897.5 | 2001.0 | 1627.3 | 2992.5 | 4060.6 |
| 2 | 91.0 | 1828.6 | 2001.0 | 1627.3 | 2860.9 | 4060.6 |
| 3 | 91.0 | 1832.4 | 2001.0 | 1627.3 | 2856.9 | 4060.6 |

### No Harden Runner
| Trial | CPU Low (m) | CPU Avg (m) | CPU High (m) | Memory Low (Mi) | Memory Avg (Mi) | Memory High (Mi) |
|-------|-------------|-------------|--------------|-----------------|------------------|------------------|
| 1 | 91.0 | 1897.5 | 2001.0 | 1627.3 | 2992.5 | 4060.6 |
| 2 | 91.0 | 1828.6 | 2001.0 | 1627.3 | 2860.9 | 4060.6 |
| 3 | 91.0 | 1832.4 | 2001.0 | 1627.3 | 2856.9 | 4060.6 |


---

## Pod Resource Utilization

### Current Harden Runner
| Trial | CPU Low (m) | CPU Avg (m) | CPU High (m) | Memory Low (Mi) | Memory Avg (Mi) | Memory High (Mi) |
|-------|-------------|-------------|--------------|-----------------|------------------|------------------|
| 1 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 |
| 2 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 |
| 3 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 |

### No Harden Runner
| Trial | CPU Low (m) | CPU Avg (m) | CPU High (m) | Memory Low (Mi) | Memory Avg (Mi) | Memory High (Mi) |
|-------|-------------|-------------|--------------|-----------------|------------------|------------------|
| 1 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 |
| 2 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 |
| 3 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 | 0.0 |


---

## Performance Comparisons

### No Harden Runner vs Current Harden Runner
| Metric | No Harden | Current | Change |
|--------|-----------|---------|--------|
| Workflow Duration | 340s | 416s | +22.35% |
