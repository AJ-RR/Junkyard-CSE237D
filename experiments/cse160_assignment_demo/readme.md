# Running `demo.py`  
A quick-start guide for launching a grading job on the phone-cluster

---

## 1.  Clone the repository
```bash
git clone https://github.com/AJ-RR/Junkyard-CSE237D.git
cd Junkyard-CSE237D/experiments/cse160_assignment_demo
```

## 2. Create/Activate a Virtual Environment
```
python3 -m venv ~/kube-env
source ~/kube-env/bin/activate
```

## 3. Ensure all dependencies are installed
```
pip install --upgrade pip
pip install kubernetes
```

## 4. Run the .py script
```
python3 demo.py
```
This should give the following output:

```bash
Job 'test-cse160-kubernetes' created successfully in namespace 'demo-may-6'.
```

## 5. ssh into the control plane
```
ssh user@10.42.0.<#phone>
```

## 6. Watch the job & view logs for the latency information
```
kubectl get pods -n demo-may-6 -o wide
```
Then using the pod name which should look like: `test-cse160-kubernetes-<#uniqueID>`
```
kubectl logs <pod-name> -n test-assignments-git
```

## 7. Delete job after finished to avoid future errors
```
kubectl delete job test-cse160-kubernetes -n demo-may-6
```
