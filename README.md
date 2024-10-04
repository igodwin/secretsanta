### Secret Santa App - Project Description

The Secret Santa app is a modern gift exchange management tool designed to streamline the organization of Secret Santa events. Initially, the app simulates the drawing of names from a pool of participants, ensuring fairness through customizable exclusion rules (e.g., no self-drawing, no drawing of significant others, etc.). As the project evolves, the goal is to scale it into a deployable service on Kubernetes with extensive features for both users and administrators.

#### Key Features:
- **Participant Management**: Add, modify, and manage participants, including specifying exclusion rules.
- **Automated Draw**: Randomly assign gift-givers to recipients with logic that respects exclusion criteria.
- **Notifications**: Notifies participants through various channels such as email, SMS, or push notifications.
  
#### Roadmap:
1. **Kubernetes Operator and CRDs**: 
   - The app will leverage a custom Kubernetes Operator to automate deployments, scaling, and updates.
   - Custom Resource Definitions (CRDs) will define the Secret Santa configuration and participants, enabling seamless integration into Kubernetes-native environments.
   
2. **Public APIs**:
   - **gRPC and REST APIs**: Expose APIs for third-party integration and front-end clients. Both gRPC and RESTful APIs will provide capabilities to create events, manage participants, run the drawing process, and query results.
   
3. **Cloud-Native and Scalable**:
   - The app will be built to run efficiently in cloud environments, with Kubernetes orchestration, auto-scaling, and load balancing.

4. **Deployment**:
   - Ready-to-deploy Helm charts and manifests will be provided for easy deployment to Kubernetes clusters.
   - The Kubernetes Operator will manage the lifecycle of the Secret Santa app, handling automated updates and configuration changes.

This project aims to serve as a comprehensive, user-friendly platform for managing Secret Santa events while demonstrating modern cloud-native application design principles.
