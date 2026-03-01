import { useMemo, useState } from 'react'

const serviceConfigs = [
  { key: 'departments', title: 'Department', endpoint: 'http://localhost:8081/departments' },
  { key: 'employees', title: 'Employee', endpoint: 'http://localhost:8082/employees' },
  { key: 'salaries', title: 'Salary', endpoint: 'http://localhost:8083/salaries' },
  { key: 'projects', title: 'Project', endpoint: 'http://localhost:8084/projects' },
]

const templates = {
  departments: { id: 'dept-001', name: 'Engineering' },
  employees: { id: 'emp-001', name: 'Alice', departmentId: 'dept-001', role: 'Backend Engineer' },
  salaries: { id: 'sal-001', employeeId: 'emp-001', monthly: 5800, currency: 'USD' },
  projects: {
    id: 'prj-001',
    name: 'ERP Core Upgrade',
    department: 'Engineering',
    members: ['emp-001'],
    timeline: 'Q2-Q4 2026',
    startDate: '2026-04-01',
    expectedEnd: '2026-12-15',
  },
}

export function App() {
  const [selected, setSelected] = useState('departments')
  const [output, setOutput] = useState('Run list/create requests from this panel.')
  const [payload, setPayload] = useState(JSON.stringify(templates.departments, null, 2))

  const active = useMemo(() => serviceConfigs.find((item) => item.key === selected), [selected])

  const runRequest = async (method) => {
    try {
      const body = method === 'POST' ? JSON.parse(payload) : undefined
      const response = await fetch(active.endpoint, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: body ? JSON.stringify(body) : undefined,
      })
      const data = await response.json()
      setOutput(JSON.stringify(data, null, 2))
    } catch (error) {
      setOutput(`Request failed: ${error.message}`)
    }
  }

  const onSelectChange = (key) => {
    setSelected(key)
    setPayload(JSON.stringify(templates[key], null, 2))
  }

  return (
    <main className="container">
      <h1>ERP Microservices Management Console</h1>
      <p>This React frontend talks to four Go services and demonstrates Kafka-based ERP architecture.</p>
      <div className="card">
        <label>
          Domain
          <select value={selected} onChange={(e) => onSelectChange(e.target.value)}>
            {serviceConfigs.map((item) => (
              <option key={item.key} value={item.key}>{item.title}</option>
            ))}
          </select>
        </label>
        <div className="buttons">
          <button onClick={() => runRequest('GET')}>List</button>
          <button onClick={() => runRequest('POST')}>Create</button>
        </div>
        <label>
          JSON Payload
          <textarea value={payload} onChange={(e) => setPayload(e.target.value)} rows={12} />
        </label>
      </div>
      <pre>{output}</pre>
    </main>
  )
}
